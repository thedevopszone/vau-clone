# Vault Clone - Quick Start Guide

## Deine aktuellen Credentials

**WICHTIG: Speichere diese Werte sicher!**

```
Root Token: EWi2lurTN2FKuyuRSQI7-SjQSuo118OEfGl3_k5gIcU=
Unseal Key: Esk4W/UeYLHJgzAXp/f2iGTpSixnrEeOpRlirUYqOJ0=
```

## Server läuft bereits im Hintergrund

Der Server läuft auf: `http://127.0.0.1:8200`

## Schnellstart

### 1. Token als Umgebungsvariable setzen
```bash
export VAULT_TOKEN=EWi2lurTN2FKuyuRSQI7-SjQSuo118OEfGl3_k5gIcU=
```

### 2. Secrets schreiben
```bash
./vault-cli write secret/database username=admin password=geheim123
./vault-cli write secret/api key=abc123 secret=xyz789
```

### 3. Secrets lesen
```bash
./vault-cli read secret/database
./vault-cli read secret/api
```

### 4. Alle Secrets auflisten
```bash
./vault-cli list
```

### 5. Secret löschen
```bash
./vault-cli delete secret/database
```

---

## Nach Server-Neustart

Wenn du den Server neu startest, führe folgende Schritte aus:

```bash
# 1. Server starten
./vault-server

# 2. Vault entsiegeln (in neuem Terminal)
./vault-cli unseal Esk4W/UeYLHJgzAXp/f2iGTpSixnrEeOpRlirUYqOJ0=

# 3. Token setzen
export VAULT_TOKEN=EWi2lurTN2FKuyuRSQI7-SjQSuo118OEfGl3_k5gIcU=

# 4. Token authentifizieren
./vault-cli auth

# 5. Jetzt kannst du Secrets nutzen
./vault-cli list
```

---

## Alle Befehle

```bash
# Status prüfen
./vault-cli status

# Vault entsiegeln
./vault-cli unseal <unseal-key>

# Vault versiegeln
./vault-cli seal

# Token authentifizieren
./vault-cli auth

# Secret schreiben
./vault-cli write <path> key1=value1 key2=value2

# Secret lesen
./vault-cli read <path>

# Secret löschen
./vault-cli delete <path>

# Alle Secrets auflisten
./vault-cli list

# Neuen Token erstellen (mit 24h Gültigkeit)
./vault-cli token-create 24h
```

---

## Beispiele

### Datenbank-Credentials speichern
```bash
./vault-cli write secret/prod/database \
  host=db.example.com \
  username=admin \
  password=super-secret-123 \
  port=5432
```

### API Keys speichern
```bash
./vault-cli write secret/prod/api \
  api_key=abc123xyz \
  api_secret=secret789 \
  endpoint=https://api.example.com
```

### Credentials lesen
```bash
./vault-cli read secret/prod/database
./vault-cli read secret/prod/api
```

---

## Server beenden

```bash
pkill -f vault-server
```

---

## Von vorne beginnen

```bash
# Server stoppen
pkill -f vault-server

# Daten löschen
rm -rf vault-data/

# Server starten
./vault-server

# Neu initialisieren
./vault-cli init
# WICHTIG: Neue Credentials notieren!
```
