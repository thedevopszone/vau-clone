package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// TokenStore manages authentication tokens
type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*Token
}

// Token represents an authentication token
type Token struct {
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time
	IsRoot    bool
}

// NewTokenStore creates a new token store
func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]*Token),
	}
}

// CreateToken creates a new token
func (ts *TokenStore) CreateToken(tokenID string, isRoot bool, ttl time.Duration) *Token {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	// Root tokens never expire (ttl = 0)
	if isRoot && ttl == 0 {
		expiresAt = time.Now().Add(100 * 365 * 24 * time.Hour) // 100 years
	}

	token := &Token{
		ID:        tokenID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsRoot:    isRoot,
	}

	ts.tokens[tokenID] = token
	return token
}

// ValidateToken checks if a token is valid
func (ts *TokenStore) ValidateToken(tokenID string) error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	token, exists := ts.tokens[tokenID]
	if !exists {
		return errors.New("invalid token")
	}

	if time.Now().After(token.ExpiresAt) {
		return errors.New("token expired")
	}

	return nil
}

// RevokeToken revokes a token
func (ts *TokenStore) RevokeToken(tokenID string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tokens[tokenID]; !exists {
		return errors.New("token not found")
	}

	delete(ts.tokens, tokenID)
	return nil
}

// IsRootToken checks if a token is a root token
func (ts *TokenStore) IsRootToken(tokenID string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	token, exists := ts.tokens[tokenID]
	if !exists {
		return false
	}

	return token.IsRoot
}

// HashToken creates a SHA-256 hash of a token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
