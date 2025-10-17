package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"vault-clone/pkg/vault"
)

var (
	vaultInstance *vault.Vault
	addr          = flag.String("addr", "127.0.0.1:8200", "HTTP server address")
	storagePath   = flag.String("storage", "./vault-data", "Storage directory path")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
}

type SecretRequest struct {
	Data map[string]interface{} `json:"data"`
}

type SecretResponse struct {
	Data map[string]interface{} `json:"data"`
}

type TokenCreateRequest struct {
	TTL string `json:"ttl"` // Duration string like "1h", "24h", etc.
}

type TokenCreateResponse struct {
	Token string `json:"token"`
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Vault-Token")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func getTokenFromHeader(r *http.Request) string {
	return r.Header.Get("X-Vault-Token")
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Status endpoint
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Initialized: vaultInstance.IsInitialized(),
		Sealed:      vaultInstance.IsSealed(),
	}
	writeJSON(w, http.StatusOK, response)
}

// Initialize endpoint
func initHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if vaultInstance.IsInitialized() {
		writeError(w, http.StatusBadRequest, "vault is already initialized")
		return
	}

	initResp, err := vaultInstance.Initialize()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, initResp)
}

// Unseal endpoint
func unsealHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Key string `json:"key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := vaultInstance.Unseal(req.Key); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, StatusResponse{
		Initialized: vaultInstance.IsInitialized(),
		Sealed:      vaultInstance.IsSealed(),
	})
}

// Seal endpoint
func sealHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := vaultInstance.Seal(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, StatusResponse{
		Initialized: vaultInstance.IsInitialized(),
		Sealed:      vaultInstance.IsSealed(),
	})
}

// Write secret endpoint
func writeSecretHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	path := r.URL.Path[len("/v1/secret/"):]
	if path == "" {
		writeError(w, http.StatusBadRequest, "missing secret path")
		return
	}

	var req SecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := vaultInstance.WriteSecret(token, path, req.Data); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// Read secret endpoint
func readSecretHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	path := r.URL.Path[len("/v1/secret/"):]
	if path == "" {
		writeError(w, http.StatusBadRequest, "missing secret path")
		return
	}

	secret, err := vaultInstance.ReadSecret(token, path)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, SecretResponse{Data: secret.Data})
}

// Delete secret endpoint
func deleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	path := r.URL.Path[len("/v1/secret/"):]
	if path == "" {
		writeError(w, http.StatusBadRequest, "missing secret path")
		return
	}

	if err := vaultInstance.DeleteSecret(token, path); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// List secrets endpoint
func listSecretsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	prefix := r.URL.Query().Get("prefix")
	secrets, err := vaultInstance.ListSecrets(token, prefix)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"keys": secrets})
}

// Create token endpoint
func createTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	var req TokenCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ttl := 24 * time.Hour // Default TTL
	if req.TTL != "" {
		parsedTTL, err := time.ParseDuration(req.TTL)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid ttl format")
			return
		}
		ttl = parsedTTL
	}

	newToken, err := vaultInstance.CreateToken(token, ttl)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, TokenCreateResponse{Token: newToken})
}

// Authenticate root token endpoint
func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := getTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing token")
		return
	}

	if err := vaultInstance.AuthenticateRootToken(token); err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "authenticated"})
}

// Router to handle secret endpoints
func secretRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		readSecretHandler(w, r)
	case http.MethodPost, http.MethodPut:
		writeSecretHandler(w, r)
	case http.MethodDelete:
		deleteSecretHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func main() {
	flag.Parse()

	// Initialize vault
	var err error
	vaultInstance, err = vault.New(*storagePath)
	if err != nil {
		log.Fatalf("Failed to create vault: %v", err)
	}

	// Setup routes with CORS middleware
	http.HandleFunc("/v1/sys/health", corsMiddleware(healthHandler))
	http.HandleFunc("/v1/sys/status", corsMiddleware(statusHandler))
	http.HandleFunc("/v1/sys/init", corsMiddleware(initHandler))
	http.HandleFunc("/v1/sys/unseal", corsMiddleware(unsealHandler))
	http.HandleFunc("/v1/sys/seal", corsMiddleware(sealHandler))
	http.HandleFunc("/v1/secret/", corsMiddleware(secretRouter))
	http.HandleFunc("/v1/secrets/list", corsMiddleware(listSecretsHandler))
	http.HandleFunc("/v1/auth/token/create", corsMiddleware(createTokenHandler))
	http.HandleFunc("/v1/auth/token/authenticate", corsMiddleware(authenticateHandler))

	fmt.Printf("Vault server starting on %s\n", *addr)
	fmt.Println("Storage path:", *storagePath)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
