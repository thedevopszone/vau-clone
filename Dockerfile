# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o vault-server ./cmd/vault-server
RUN CGO_ENABLED=0 GOOS=linux go build -o vault-cli ./cmd/vault-cli

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/vault-server .
COPY --from=builder /app/vault-cli .

# Create data directory
RUN mkdir -p /vault-data

# Expose port
EXPOSE 8200

# Run server
CMD ["./vault-server", "-addr", "0.0.0.0:8200", "-storage", "/vault-data"]
