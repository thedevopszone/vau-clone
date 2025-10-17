#!/bin/bash

# Helper script to interact with Vault running in Docker

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if docker-compose is running
if ! docker-compose ps | grep -q "vault.*Up"; then
    echo -e "${RED}Error: Vault container is not running!${NC}"
    echo "Start it with: docker-compose up -d"
    exit 1
fi

# If VAULT_TOKEN is set, use it
if [ -n "$VAULT_TOKEN" ]; then
    docker-compose exec vault sh -c "export VAULT_TOKEN=$VAULT_TOKEN && ./vault-cli $*"
else
    # For commands that don't need token (init, status, unseal)
    docker-compose exec vault ./vault-cli "$@"
fi
