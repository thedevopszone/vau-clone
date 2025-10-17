# Vault Clone

A simplified HashiCorp Vault clone written in Go with basic secret management functionality.

## Features

- **Encryption**: AES-256-GCM encryption for all secrets
- **Storage**: File-based storage backend with JSON persistence
- **Authentication**: Token-based authentication system
- **Seal/Unseal**: Vault can be sealed and unsealed with a master key
- **HTTP API**: RESTful API for all operations
- **CLI Client**: User-friendly command-line interface

## Architecture

```
vault-clone/
├── cmd/
│   ├── vault-server/    # HTTP API server
│   └── vault-cli/       # CLI client
├── pkg/
│   ├── auth/           # Authentication and token management
│   ├── crypto/         # Encryption/decryption operations
│   ├── storage/        # Storage backend interface
│   └── vault/          # Core vault logic
└── vault-data/         # Storage directory (created at runtime)
```

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:
```bash
cd vault-clone
```

2. Install dependencies:
```bash
go mod download
```

3. Build the binaries:
```bash
go build -o vault-server ./cmd/vault-server
go build -o vault-cli ./cmd/vault-cli
```

### Running the Vault Server

Start the server (default address: 127.0.0.1:8200):
```bash
./vault-server
```

Or specify a custom address and storage path:
```bash
./vault-server -addr 127.0.0.1:8300 -storage ./my-vault-data
```

### Using the CLI

Set the vault address (if not using default):
```bash
export VAULT_ADDR=http://127.0.0.1:8200
```

#### 1. Initialize the Vault

```bash
./vault-cli init
```

This will output:
- **Root Token**: Used for authentication
- **Unseal Key**: Used to unseal the vault

**IMPORTANT**: Save these credentials securely! They are only shown once.

#### 2. Unseal the Vault

```bash
./vault-cli unseal <unseal-key>
```

#### 3. Set Authentication Token

```bash
export VAULT_TOKEN=<root-token>
```

#### 4. Write Secrets

```bash
./vault-cli write secret/myapp password=secret123 api_key=abc123
```

#### 5. Read Secrets

```bash
./vault-cli read secret/myapp
```

#### 6. List Secrets

```bash
./vault-cli list
```

List with prefix:
```bash
./vault-cli list secret/
```

#### 7. Delete Secrets

```bash
./vault-cli delete secret/myapp
```

#### 8. Create New Tokens

```bash
./vault-cli token-create 24h
```

#### 9. Seal the Vault

```bash
./vault-cli seal
```

## API Endpoints

### System Operations

- `GET /v1/sys/health` - Health check
- `GET /v1/sys/status` - Get vault status (initialized, sealed)
- `POST /v1/sys/init` - Initialize the vault
- `POST /v1/sys/unseal` - Unseal the vault
- `POST /v1/sys/seal` - Seal the vault

### Secret Operations

- `POST /v1/secret/:path` - Write a secret
- `GET /v1/secret/:path` - Read a secret
- `DELETE /v1/secret/:path` - Delete a secret
- `GET /v1/secrets/list?prefix=` - List secrets

### Authentication

- `POST /v1/auth/token/create` - Create a new token

## Example Usage

### Complete Workflow

```bash
# Terminal 1: Start the server
./vault-server

# Terminal 2: Initialize and use the vault
./vault-cli init
# Save the output credentials!

./vault-cli unseal <unseal-key>
export VAULT_TOKEN=<root-token>

# Store secrets
./vault-cli write secret/database username=admin password=secretpass
./vault-cli write secret/api key=abc123 endpoint=https://api.example.com

# Read secrets
./vault-cli read secret/database
./vault-cli read secret/api

# List all secrets
./vault-cli list

# Create a new token with 1 hour TTL
./vault-cli token-create 1h

# Seal the vault when done
./vault-cli seal
```

### Using the API Directly

```bash
# Initialize
curl -X POST http://127.0.0.1:8200/v1/sys/init

# Unseal
curl -X POST http://127.0.0.1:8200/v1/sys/unseal \
  -H "Content-Type: application/json" \
  -d '{"key":"<unseal-key>"}'

# Write secret
curl -X POST http://127.0.0.1:8200/v1/secret/myapp \
  -H "X-Vault-Token: <token>" \
  -H "Content-Type: application/json" \
  -d '{"data":{"password":"secret123"}}'

# Read secret
curl -X GET http://127.0.0.1:8200/v1/secret/myapp \
  -H "X-Vault-Token: <token>"
```

## Security Features

- **Encryption at Rest**: All secrets are encrypted with AES-256-GCM before storage
- **Token-based Authentication**: All operations require valid authentication tokens
- **Seal/Unseal Mechanism**: Vault must be unsealed to access secrets
- **Key Derivation**: PBKDF2 for secure key derivation from passwords
- **Secure Token Generation**: Cryptographically secure random token generation

## Limitations

This is a simplified clone for educational purposes. It lacks many features of production Vault:

- No TLS/HTTPS support (use a reverse proxy in production)
- Single unseal key (production Vault uses Shamir's Secret Sharing)
- File-based storage only (no distributed backends)
- Limited authentication methods (token-only)
- No audit logging
- No policy system
- No secret rotation
- No high availability

## Environment Variables

- `VAULT_ADDR` - Vault server address (default: http://127.0.0.1:8200)
- `VAULT_TOKEN` - Authentication token for CLI operations

## Contributing

This is an educational project demonstrating basic secret management concepts. Feel free to extend it with additional features!

## License

MIT License

## Acknowledgments

Inspired by HashiCorp Vault's architecture and functionality.
