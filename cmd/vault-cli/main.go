package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	defaultAddr = "http://127.0.0.1:8200"
)

type InitResponse struct {
	RootToken string `json:"root_token"`
	UnsealKey string `json:"unseal_key"`
}

type StatusResponse struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
}

type SecretResponse struct {
	Data map[string]interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func getVaultAddr() string {
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		return addr
	}
	return defaultAddr
}

func getVaultToken() string {
	return os.Getenv("VAULT_TOKEN")
}

func makeRequest(method, endpoint string, body interface{}, token string) (*http.Response, error) {
	addr := getVaultAddr()
	url := addr + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("X-Vault-Token", token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func printUsage() {
	fmt.Println("Vault CLI - A simple HashiCorp Vault clone")
	fmt.Println("\nUsage:")
	fmt.Println("  vault-cli <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  status                           Show vault status")
	fmt.Println("  init                             Initialize the vault")
	fmt.Println("  unseal <key>                     Unseal the vault")
	fmt.Println("  seal                             Seal the vault")
	fmt.Println("  auth                             Authenticate root token")
	fmt.Println("  write <path> <key=value>...      Write a secret")
	fmt.Println("  read <path>                      Read a secret")
	fmt.Println("  delete <path>                    Delete a secret")
	fmt.Println("  list [prefix]                    List secrets")
	fmt.Println("  token-create [ttl]               Create a new token")
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  VAULT_ADDR      Vault server address (default: http://127.0.0.1:8200)")
	fmt.Println("  VAULT_TOKEN     Authentication token")
	fmt.Println("\nExamples:")
	fmt.Println("  vault-cli init")
	fmt.Println("  vault-cli unseal <unseal-key>")
	fmt.Println("  export VAULT_TOKEN=<root-token>")
	fmt.Println("  vault-cli auth")
	fmt.Println("  vault-cli write secret/myapp password=secret123 api_key=abc123")
	fmt.Println("  vault-cli read secret/myapp")
	fmt.Println("  vault-cli delete secret/myapp")
	fmt.Println("  vault-cli list")
}

func handleStatus() error {
	resp, err := makeRequest("GET", "/v1/sys/status", nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return err
	}

	fmt.Printf("Initialized: %v\n", status.Initialized)
	fmt.Printf("Sealed: %v\n", status.Sealed)
	return nil
}

func handleInit() error {
	resp, err := makeRequest("POST", "/v1/sys/init", nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("init failed: %s", errResp.Error)
	}

	var initResp InitResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return err
	}

	fmt.Println("Vault initialized successfully!")
	fmt.Println("\nIMPORTANT: Save these credentials securely!")
	fmt.Printf("\nRoot Token: %s\n", initResp.RootToken)
	fmt.Printf("Unseal Key: %s\n", initResp.UnsealKey)
	fmt.Println("\nTo unseal the vault, run:")
	fmt.Printf("  vault-cli unseal %s\n", initResp.UnsealKey)
	fmt.Println("\nTo authenticate, set the token:")
	fmt.Printf("  export VAULT_TOKEN=%s\n", initResp.RootToken)
	return nil
}

func handleUnseal(key string) error {
	body := map[string]string{"key": key}
	resp, err := makeRequest("POST", "/v1/sys/unseal", body, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("unseal failed: %s", errResp.Error)
	}

	fmt.Println("Vault unsealed successfully!")
	return handleStatus()
}

func handleSeal() error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	resp, err := makeRequest("POST", "/v1/sys/seal", nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("seal failed: %s", errResp.Error)
	}

	fmt.Println("Vault sealed successfully!")
	return nil
}

func handleWrite(path string, kvPairs []string) error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	data := make(map[string]interface{})
	for _, pair := range kvPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid key=value pair: %s", pair)
		}
		data[parts[0]] = parts[1]
	}

	body := map[string]interface{}{"data": data}
	resp, err := makeRequest("POST", "/v1/secret/"+path, body, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("write failed: %s", errResp.Error)
	}

	fmt.Printf("Secret written successfully to: %s\n", path)
	return nil
}

func handleRead(path string) error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	resp, err := makeRequest("GET", "/v1/secret/"+path, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("read failed: %s", errResp.Error)
	}

	var secretResp SecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&secretResp); err != nil {
		return err
	}

	fmt.Printf("Secret at %s:\n", path)
	for key, value := range secretResp.Data {
		fmt.Printf("  %s: %v\n", key, value)
	}
	return nil
}

func handleDelete(path string) error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	resp, err := makeRequest("DELETE", "/v1/secret/"+path, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("delete failed: %s", errResp.Error)
	}

	fmt.Printf("Secret deleted successfully: %s\n", path)
	return nil
}

func handleList(prefix string) error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	endpoint := "/v1/secrets/list"
	if prefix != "" {
		endpoint += "?prefix=" + prefix
	}

	resp, err := makeRequest("GET", endpoint, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("list failed: %s", errResp.Error)
	}

	var listResp struct {
		Keys []string `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return err
	}

	if len(listResp.Keys) == 0 {
		fmt.Println("No secrets found")
		return nil
	}

	fmt.Println("Secrets:")
	for _, key := range listResp.Keys {
		fmt.Printf("  %s\n", key)
	}
	return nil
}

func handleTokenCreate(ttl string) error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	body := make(map[string]string)
	if ttl != "" {
		body["ttl"] = ttl
	}

	resp, err := makeRequest("POST", "/v1/auth/token/create", body, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("token create failed: %s", errResp.Error)
	}

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	fmt.Printf("New token created: %s\n", tokenResp.Token)
	return nil
}

func handleAuth() error {
	token := getVaultToken()
	if token == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	resp, err := makeRequest("POST", "/v1/auth/token/authenticate", nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("authentication failed: %s", errResp.Error)
	}

	fmt.Println("Token authenticated successfully!")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	var err error

	switch command {
	case "status":
		err = handleStatus()
	case "init":
		err = handleInit()
	case "unseal":
		if len(os.Args) < 3 {
			fmt.Println("Error: unseal key required")
			os.Exit(1)
		}
		err = handleUnseal(os.Args[2])
	case "seal":
		err = handleSeal()
	case "auth":
		err = handleAuth()
	case "write":
		if len(os.Args) < 4 {
			fmt.Println("Error: path and at least one key=value pair required")
			os.Exit(1)
		}
		err = handleWrite(os.Args[2], os.Args[3:])
	case "read":
		if len(os.Args) < 3 {
			fmt.Println("Error: path required")
			os.Exit(1)
		}
		err = handleRead(os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: path required")
			os.Exit(1)
		}
		err = handleDelete(os.Args[2])
	case "list":
		prefix := ""
		if len(os.Args) >= 3 {
			prefix = os.Args[2]
		}
		err = handleList(prefix)
	case "token-create":
		ttl := ""
		if len(os.Args) >= 3 {
			ttl = os.Args[2]
		}
		err = handleTokenCreate(ttl)
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
