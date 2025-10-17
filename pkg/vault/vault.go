package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"vault-clone/pkg/auth"
	"vault-clone/pkg/crypto"
	"vault-clone/pkg/storage"
)

// Vault represents the main vault instance
type Vault struct {
	storage      storage.Storage
	tokenStore   *auth.TokenStore
	mu           sync.RWMutex
	sealed       bool
	initialized  bool
	encryptionKey []byte
	rootToken    string
}

// Secret represents a secret stored in the vault
type Secret struct {
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// InitResponse contains the initialization response
type InitResponse struct {
	RootToken     string `json:"root_token"`
	UnsealKey     string `json:"unseal_key"`
}

// New creates a new vault instance
func New(storagePath string) (*Vault, error) {
	store, err := storage.NewFileStorage(storagePath)
	if err != nil {
		return nil, err
	}

	v := &Vault{
		storage:    store,
		tokenStore: auth.NewTokenStore(),
		sealed:     true,
		initialized: false,
	}

	// Check if vault is already initialized
	if err := v.checkInitialized(); err == nil {
		v.initialized = true
	}

	return v, nil
}

// Initialize initializes the vault and returns the root token and unseal key
func (v *Vault) Initialize() (*InitResponse, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.initialized {
		return nil, errors.New("vault is already initialized")
	}

	// Generate unseal key (master key)
	unsealKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	// Generate root token
	rootTokenRaw, err := crypto.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Store the unseal key securely (encrypted with itself for verification)
	encryptedKey, err := crypto.Encrypt(unsealKey, unsealKey)
	if err != nil {
		return nil, err
	}

	if err := v.storage.Put("core/unseal-key", []byte(encryptedKey)); err != nil {
		return nil, err
	}

	// Store root token hash
	rootTokenHash := auth.HashToken(rootTokenRaw)
	if err := v.storage.Put("core/root-token", []byte(rootTokenHash)); err != nil {
		return nil, err
	}

	// Create root token in token store
	v.tokenStore.CreateToken(rootTokenRaw, true, 0) // Root token never expires

	v.initialized = true
	v.rootToken = rootTokenRaw

	return &InitResponse{
		RootToken: rootTokenRaw,
		UnsealKey: base64.StdEncoding.EncodeToString(unsealKey),
	}, nil
}

// Unseal unseals the vault with the unseal key
func (v *Vault) Unseal(unsealKeyStr string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return errors.New("vault is not initialized")
	}

	if !v.sealed {
		return errors.New("vault is already unsealed")
	}

	unsealKey, err := base64.StdEncoding.DecodeString(unsealKeyStr)
	if err != nil {
		return errors.New("invalid unseal key format")
	}

	// Verify unseal key
	encryptedKeyData, err := v.storage.Get("core/unseal-key")
	if err != nil {
		return err
	}

	_, err = crypto.Decrypt(string(encryptedKeyData), unsealKey)
	if err != nil {
		return errors.New("invalid unseal key")
	}

	v.encryptionKey = unsealKey
	v.sealed = false

	// Restore root token to token store after unseal
	rootTokenHashData, err := v.storage.Get("core/root-token")
	if err == nil && len(rootTokenHashData) > 0 {
		// Note: We can't recreate the actual token from its hash
		// Root token must be provided via storage or environment
		// For now, we'll rely on users to keep their root token
	}

	return nil
}

// Seal seals the vault
func (v *Vault) Seal() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.sealed {
		return errors.New("vault is already sealed")
	}

	v.encryptionKey = nil
	v.sealed = true

	return nil
}

// IsSealed returns whether the vault is sealed
func (v *Vault) IsSealed() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.sealed
}

// IsInitialized returns whether the vault is initialized
func (v *Vault) IsInitialized() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.initialized
}

// WriteSecret writes a secret to the vault
func (v *Vault) WriteSecret(token, path string, data map[string]interface{}) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.sealed {
		return errors.New("vault is sealed")
	}

	if err := v.tokenStore.ValidateToken(token); err != nil {
		return err
	}

	secret := &Secret{
		Data:      data,
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	secretJSON, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	// Encrypt the secret
	encrypted, err := crypto.Encrypt(secretJSON, v.encryptionKey)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("secret/%s", path)
	return v.storage.Put(key, []byte(encrypted))
}

// ReadSecret reads a secret from the vault
func (v *Vault) ReadSecret(token, path string) (*Secret, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.sealed {
		return nil, errors.New("vault is sealed")
	}

	if err := v.tokenStore.ValidateToken(token); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("secret/%s", path)
	encryptedData, err := v.storage.Get(key)
	if err != nil {
		return nil, err
	}

	// Decrypt the secret
	decrypted, err := crypto.Decrypt(string(encryptedData), v.encryptionKey)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := json.Unmarshal(decrypted, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

// DeleteSecret deletes a secret from the vault
func (v *Vault) DeleteSecret(token, path string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.sealed {
		return errors.New("vault is sealed")
	}

	if err := v.tokenStore.ValidateToken(token); err != nil {
		return err
	}

	key := fmt.Sprintf("secret/%s", path)
	return v.storage.Delete(key)
}

// ListSecrets lists all secrets with the given prefix
func (v *Vault) ListSecrets(token, prefix string) ([]string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.sealed {
		return nil, errors.New("vault is sealed")
	}

	if err := v.tokenStore.ValidateToken(token); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("secret/%s", prefix)
	keys, err := v.storage.List(key)
	if err != nil {
		return nil, err
	}

	// Strip the "secret/" prefix from keys
	var secrets []string
	for _, k := range keys {
		if len(k) > 7 {
			secrets = append(secrets, k[7:])
		}
	}

	return secrets, nil
}

// CreateToken creates a new authentication token
func (v *Vault) CreateToken(rootToken string, ttl time.Duration) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.sealed {
		return "", errors.New("vault is sealed")
	}

	if err := v.tokenStore.ValidateToken(rootToken); err != nil {
		return "", err
	}

	if !v.tokenStore.IsRootToken(rootToken) {
		return "", errors.New("only root token can create new tokens")
	}

	newToken, err := crypto.GenerateToken()
	if err != nil {
		return "", err
	}

	v.tokenStore.CreateToken(newToken, false, ttl)
	return newToken, nil
}

// checkInitialized checks if the vault has been initialized
func (v *Vault) checkInitialized() error {
	_, err := v.storage.Get("core/unseal-key")
	return err
}

// GetRootToken returns the root token (only available immediately after initialization)
func (v *Vault) GetRootToken() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.rootToken
}

// AuthenticateRootToken verifies and registers a root token in the token store
func (v *Vault) AuthenticateRootToken(token string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.sealed {
		return errors.New("vault is sealed")
	}

	// Get stored root token hash
	rootTokenHashData, err := v.storage.Get("core/root-token")
	if err != nil {
		return errors.New("no root token configured")
	}

	// Verify the provided token matches the stored hash
	providedHash := auth.HashToken(token)
	if providedHash != string(rootTokenHashData) {
		return errors.New("invalid root token")
	}

	// Add token to token store with no expiration
	v.tokenStore.CreateToken(token, true, 0)
	return nil
}
