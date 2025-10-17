# Docker Setup für Vault Clone

## Schnellstart mit Docker Compose

### 1. Container starten
```bash
docker-compose up -d
```

### 2. Logs ansehen
```bash
docker-compose logs -f vault
```

### 3. Vault initialisieren
```bash
# CLI im Container ausführen
docker-compose exec vault ./vault-cli init
```

**WICHTIG: Speichere Root Token und Unseal Key!**

### 4. Vault entsiegeln
```bash
docker-compose exec vault ./vault-cli unseal <dein-unseal-key>
```

### 5. Token authentifizieren
```bash
docker-compose exec vault sh -c "export VAULT_TOKEN=<dein-root-token> && ./vault-cli auth"
```

### 6. Secrets verwalten
```bash
# Secret schreiben
docker-compose exec vault sh -c "export VAULT_TOKEN=<dein-root-token> && ./vault-cli write secret/database username=admin password=secret123"

# Secret lesen
docker-compose exec vault sh -c "export VAULT_TOKEN=<dein-root-token> && ./vault-cli read secret/database"

# Alle Secrets auflisten
docker-compose exec vault sh -c "export VAULT_TOKEN=<dein-root-token> && ./vault-cli list"
```

---

## Einfacherer Workflow mit Umgebungsvariable

### Erstelle ein Helper-Script
```bash
cat > vault-cli-docker.sh << 'EOF'
#!/bin/bash
docker-compose exec vault sh -c "export VAULT_TOKEN=${VAULT_TOKEN} && ./vault-cli $*"
EOF

chmod +x vault-cli-docker.sh
```

### Dann kannst du es so nutzen:
```bash
# Token setzen
export VAULT_TOKEN=<dein-root-token>

# Befehle ausführen
./vault-cli-docker.sh status
./vault-cli-docker.sh unseal <unseal-key>
./vault-cli-docker.sh auth
./vault-cli-docker.sh write secret/test key=value
./vault-cli-docker.sh read secret/test
./vault-cli-docker.sh list
```

---

## Von Host-System auf Vault zugreifen

Der Vault ist auf Port 8200 verfügbar. Du kannst auch den lokalen CLI nutzen:

```bash
# Token setzen
export VAULT_TOKEN=<dein-root-token>
export VAULT_ADDR=http://localhost:8200

# Lokalen CLI verwenden
./vault-cli status
./vault-cli auth
./vault-cli write secret/myapp password=secret123
./vault-cli read secret/myapp
```

---

## Nützliche Docker Compose Befehle

```bash
# Container starten
docker-compose up -d

# Container stoppen
docker-compose stop

# Container stoppen und entfernen
docker-compose down

# Container stoppen und Daten löschen
docker-compose down -v

# Logs ansehen
docker-compose logs -f vault

# Container Status prüfen
docker-compose ps

# In Container shell einsteigen
docker-compose exec vault sh

# Container neu bauen
docker-compose build --no-cache
docker-compose up -d
```

---

## Daten-Persistenz

Die Vault-Daten werden in einem Docker Volume gespeichert:
- Volume Name: `vault-clone_vault-data`
- Gemountet in: `/vault-data` im Container

### Volume verwalten
```bash
# Volumes anzeigen
docker volume ls

# Volume inspizieren
docker volume inspect vault-clone_vault-data

# Volume löschen (nur wenn Container gestoppt)
docker-compose down -v
```

---

## Produktions-Setup (mit docker-compose.prod.yml)

Für Produktion kannst du zusätzliche Konfigurationen hinzufügen:

```yaml
version: '3.8'

services:
  vault:
    build: .
    container_name: vault-clone-prod
    ports:
      - "8200:8200"
    volumes:
      - ./vault-data-prod:/vault-data
    environment:
      - VAULT_ADDR=http://0.0.0.0:8200
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

Starten mit:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

---

## Troubleshooting

### Container startet nicht
```bash
# Logs prüfen
docker-compose logs vault

# Container Status
docker-compose ps
```

### Port bereits belegt
```bash
# Anderen Port verwenden (z.B. 8201)
# In docker-compose.yml ändern:
ports:
  - "8201:8200"

# Dann VAULT_ADDR anpassen:
export VAULT_ADDR=http://localhost:8201
```

### Daten zurücksetzen
```bash
# Alles löschen und neu starten
docker-compose down -v
docker-compose up -d
docker-compose exec vault ./vault-cli init
```

---

## Mit curl testen

```bash
# Status prüfen
curl http://localhost:8200/v1/sys/status

# Secret schreiben
curl -X POST http://localhost:8200/v1/secret/test \
  -H "X-Vault-Token: <dein-token>" \
  -H "Content-Type: application/json" \
  -d '{"data":{"password":"secret123"}}'

# Secret lesen
curl -X GET http://localhost:8200/v1/secret/test \
  -H "X-Vault-Token: <dein-token>"
```
