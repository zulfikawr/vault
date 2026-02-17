# Quick Start Guide

Get up and running with Vault in 5 minutes.

## Prerequisites

- Go 1.25.6 or later (for building from source)
- Or download a pre-built binary

## Installation

### Option 1: One-line Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/zulfikawr/vault/main/install.sh | bash
```

### Option 2: Download Binary

Download from [GitHub Releases](https://github.com/zulfikawr/vault/releases) and add to your PATH.

### Option 3: Build from Source

```bash
git clone https://github.com/zulfikawr/vault.git
cd vault
go build -o vault ./cmd/vault
sudo mv vault /usr/local/bin/
```

## Step 1: Initialize Vault

Create a new admin user and initialize the database:

```bash
vault init --email "admin@example.com" --username "admin" --password "yourpassword"
```

**Output:**
```
🚀 Initializing Vault project...

✓ Created data directory: ./vault_data
✓ Generated config.json
✓ Created .env.example
✓ Initialized database: vault_data/vault.db
✓ Created system collections

👤 Create your first admin user:
   Email: admin@example.com
   Username: admin
   Password: ********

✓ Admin user created successfully

🎉 Vault initialized successfully!

Next steps:
  Run: vault serve
  Visit: http://localhost:8090/
```

## Step 2: Start the Server

```bash
vault serve
```

**Output:**
```
Starting Vault server on port 8090...
Database initialized at vault_data/vault.db
Admin dashboard available at http://localhost:8090/_/
```

## Step 3: Access the Admin Dashboard

Open your browser and visit: **http://localhost:8090/_/**

Login with your admin credentials.

## Step 4: Create Your First Collection

Collections are like database tables but can be created dynamically.

### Via Admin Dashboard

1. Click "Collections" in the sidebar
2. Click "New Collection"
3. Enter collection name (e.g., `posts`)
4. Add fields:
   - `title` (text)
   - `content` (text)
   - `published` (boolean)
   - `author` (text)
5. Click "Create Collection"

### Via CLI

```bash
vault collection create \
  --name "posts" \
  --fields "title:text,content:text,published:boolean,author:text" \
  --email "admin@example.com" \
  --password "yourpassword"
```

## Step 5: Add Some Data

### Via Admin Dashboard

1. Click on your collection name
2. Click "New Record"
3. Fill in the fields
4. Click "Save"

### Via API

```bash
# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8090/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' \
  | jq -r '.token')

# Create a record
curl -X POST http://localhost:8090/api/collections/posts/records \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hello World",
    "content": "This is my first post!",
    "published": true,
    "author": "admin"
  }'
```

## Step 6: Query Your Data

```bash
# List all records
curl http://localhost:8090/api/collections/posts/records

# Filter records
curl "http://localhost:8090/api/collections/posts/records?filter=published=true"

# Get single record
curl http://localhost:8090/api/collections/posts/records/RECORD_ID
```

## What's Next?

- [Learn about Collections](./concepts/collections.md)
- [Set up Authorization Rules](./concepts/rules.md)
- [Configure File Storage](./concepts/storage.md)
- [Enable Real-time Updates](./api/realtime.md)
- [Deploy to Production](./advanced/deployment.md)

## Common Commands

```bash
# List all collections
vault collection list --email "admin@example.com" --password "yourpassword"

# Get collection details
vault collection get --name "posts" --email "admin@example.com" --password "yourpassword"

# Create backup
vault backup create --output "backup.zip"

# Export data
vault export json --output "backup.json"

# View help
vault -h
vault collection -h
```

## Default Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| Port | 8090 | HTTP server port |
| Data Directory | ./vault_data | Database and storage location |
| Database | vault_data/vault.db | SQLite database file |
| Admin Dashboard | /_/ | Dashboard URL path |
| API Endpoint | /api/ | REST API URL path |

## Troubleshooting

### Port Already in Use

```bash
vault serve --port 8091
```

### Reset Admin Password

```bash
vault admin reset-password \
  --email "admin@example.com" \
  --password "newpassword"
```

### View Logs

```bash
cat vault_data/vault.log
```

### Delete and Start Over

```bash
rm -rf vault_data config.json
vault init --email "admin@example.com" --username "admin" --password "newpassword"
```

## Need Help?

- Check the [FAQ](./troubleshooting/faq.md)
- Read the [full documentation](./README.md)
- Open an [issue on GitHub](https://github.com/zulfikawr/vault/issues)
