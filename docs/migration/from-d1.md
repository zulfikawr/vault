# Migrate from Cloudflare D1 to Vault

This guide walks you through migrating your Cloudflare D1 database and R2 storage to Vault.

## Overview

| Source | Destination | Method |
|--------|-------------|--------|
| D1 Database | Vault SQLite | SQL export/import |
| R2 Bucket | Vault Storage | File copy |

## Prerequisites

- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/install-and-update/) installed
- Access to your Cloudflare account
- Vault installed and initialized

## Step 1: Export D1 Database

### Export Schema and Data

```bash
# Export your D1 database
npx wrangler d1 export <DATABASE_NAME> --output d1-export.sql --remote
```

**Example:**
```bash
npx wrangler d1 export homepage-db --output d1-export.sql --remote
```

**Output:**
```
⛅️ wrangler 4.66.0
───────────────────
Resource location: remote 

🌀 Executing on remote database homepage-db:
├ Creating export
├ Downloading SQL to d1-export.sql
🌀 Downloaded to d1-export.sql successfully!
```

### Verify Export

```bash
# Check file size
ls -lh d1-export.sql

# Preview content
head -50 d1-export.sql
```

**Expected content:**
```sql
PRAGMA defer_foreign_keys=TRUE;
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE,
  password_hash TEXT,
  role TEXT DEFAULT 'user',
  created_at INTEGER DEFAULT (unixepoch())
);
INSERT INTO "users" VALUES('admin','zulfikawr@gmail.com','...','admin',1769936204);
```

## Step 2: Export R2 Bucket

### Option A: Using Wrangler (Recommended)

Create a script to download all objects:

```javascript
// download-r2.js
const { execSync } = require("child_process");
const { mkdirSync } = require("fs");

const BUCKET = "your-bucket-name";
const OUTPUT_DIR = "./r2-export";

mkdirSync(OUTPUT_DIR, { recursive: true });

// List objects
const list = execSync(
  `npx wrangler r2 bucket info ${BUCKET} --json`,
  { encoding: 'utf8' }
);
const info = JSON.parse(list);
console.log(`Found ${info.object_count} objects (${info.bucket_size})`);

// Download each object
const objects = execSync(
  `npx wrangler r2 object list ${BUCKET} --json`,
  { encoding: 'utf8' }
);

// Parse and download (implement based on your structure)
```

### Option B: Using rclone

```bash
# Install rclone
curl https://rclone.org/install.sh | sudo bash

# Configure R2
rclone config
# Choose "S3 Compatible"
# Provider: Cloudflare R2
# Access Key ID: your-access-key
# Secret Access Key: your-secret-key
# Endpoint: your-account-id.r2.cloudflarestorage.com

# Download bucket
rclone copy r2:your-bucket ./r2-export/
```

### Option C: Manual Download

For small buckets:

```bash
# List objects
npx wrangler r2 bucket info your-bucket

# Download individual objects
npx wrangler r2 object get "your-bucket/path/to/file" \
  --file "./r2-export/path/to/file" --remote
```

## Step 3: Initialize Vault

```bash
# Create new directory
mkdir vault-project && cd vault-project

# Initialize Vault
vault init --email "admin@example.com" \
  --username "admin" \
  --password "your-secure-password"
```

**Output:**
```
🚀 Initializing Vault project...

✓ Created data directory: ./vault_data
✓ Generated config.json
✓ Created .env.example
✓ Initialized database: vault_data/vault.db
✓ Created system collections

✓ Admin user created successfully

🎉 Vault initialized successfully!
```

## Step 4: Import D1 Data

### Dry Run First

```bash
vault import d1 d1-export.sql --dry-run
```

**Expected output:**
```
📥 Importing D1 dump: d1-export.sql
Found 3550 SQL statements

✓ Table exists: users
✓ Created table: posts
✓ Created table: projects
  ... 500 records imported
  ... 1000 records imported

✅ D1 Import complete!
   Tables: 17
   Records: 3528
   Errors: 0
```

### Import for Real

```bash
vault import d1 d1-export.sql
```

**What happens:**
1. Creates tables from D1 SQL
2. Imports all records
3. Registers collections in Vault
4. Maps field types automatically

### Verify Import

```bash
# List collections
vault collection list \
  --email "admin@example.com" \
  --password "your-secure-password"
```

**Expected output:**
```
Collections:
┌────────────────────────┬──────┬────────┐
│ Name                   │ Type │ Fields │
├────────────────────────┼──────┼────────┤
│ users                  │ base │      4 │
│ posts                  │ base │      9 │
│ projects               │ base │     13 │
│ certificates           │ base │      9 │
│ ...                    │ ...  │    ... │
└────────────────────────┴──────┴────────┘

Total: 17 collections
```

## Step 5: Import R2 Storage

```bash
# Copy R2 files to Vault storage
cp -r r2-export/* vault_data/storage/

# Verify
ls -la vault_data/storage/
```

**Preserved structure:**
```
vault_data/storage/
├── certificates/
│   ├── cs50x-2024-computer-science/
│   │   ├── image.png
│   │   └── logo.png
│   └── ...
├── employments/
│   └── ...
└── incremental-cache/
    └── ...
```

## Step 6: Test Your Migration

### Check Record Counts

```bash
# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8090/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"your-secure-password"}' \
  | jq -r '.token')

# Check record counts
curl -s "http://localhost:8090/api/collections/posts/records?page=1&perPage=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.totalRecords'

curl -s "http://localhost:8090/api/collections/projects/records?page=1&perPage=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.totalRecords'
```

### Test File Access

```bash
# List storage
vault storage list \
  --path "/" \
  --email "admin@example.com" \
  --password "your-secure-password"
```

### Verify Data Integrity

Compare record counts between D1 and Vault:

```bash
# D1 count
npx wrangler d1 execute homepage-db --remote \
  --command="SELECT COUNT(*) FROM posts"

# Vault count (should match)
curl -s "http://localhost:8090/api/collections/posts/records?page=1&perPage=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.totalRecords'
```

## Step 7: Update Your Application

### Update API Endpoints

**Before (Workers):**
```javascript
const response = await env.DB.prepare(
  "SELECT * FROM posts WHERE published = ?"
).bind(1).all();
```

**After (Vault API):**
```javascript
const response = await fetch(
  "http://localhost:8090/api/collections/posts/records?filter=published=1",
  {
    headers: {
      "Authorization": `Bearer ${token}`
    }
  }
);
const data = await response.json();
```

### Update File Paths

**Before (R2):**
```javascript
const url = `https://r2.yourdomain.com/certificates/image.png`;
```

**After (Vault Storage):**
```javascript
const url = `http://localhost:8090/api/storage/certificates/image.png`;
```

## Troubleshooting

### "Table already exists" Error

Tables are created but not registered. Run:

```bash
vault migrate sync
```

### Missing Collections in Status

```bash
# Collections exist but not showing in migrate status
vault collection list --email "admin@example.com" --password "yourpassword"
```

### Duplicate Key Errors

Some records may have conflicting IDs. Re-import with OR REPLACE:

```bash
# Edit SQL file to use INSERT OR REPLACE
sed -i 's/INSERT INTO/INSERT OR REPLACE INTO/g' d1-export.sql

# Re-import
vault import d1 d1-export.sql
```

### Large Import Timeout

For very large databases (>10,000 records):

```bash
# Increase timeout in import.go or import in batches
# Split SQL file
split -l 1000 d1-export.sql d1-batch-

# Import each batch
for file in d1-batch-*; do
  vault import d1 $file
done
```

## Post-Migration Checklist

- [ ] All tables imported successfully
- [ ] Record counts match source
- [ ] Collections registered in Vault
- [ ] File storage copied
- [ ] Authentication working
- [ ] API endpoints responding
- [ ] Admin dashboard accessible
- [ ] Authorization rules configured
- [ ] Application updated to use Vault API
- [ ] Backups configured

## Performance Tips

### Optimize After Import

```bash
# Run VACUUM to optimize database
sqlite3 vault_data/vault.db "VACUUM;"

# Analyze for query optimization
sqlite3 vault_data/vault.db "ANALYZE;"
```

### Configure for Production

```bash
# Edit config.json
{
  "port": 8090,
  "db_path": "./vault_data/vault.db",
  "data_dir": "./vault_data",
  "log_level": "INFO",
  "rate_limit_per_min": 300,
  "max_file_upload_size": 10485760
}
```

## Rollback Plan

If you need to rollback to D1:

1. Keep D1 database intact until fully validated
2. Export Vault data before any changes:
   ```bash
   vault export sql --output vault-backup.sql
   ```
3. Update application config to point back to D1

## Next Steps

- [Configure Authorization Rules](../concepts/rules.md)
- [Set up Real-time Subscriptions](../api/realtime.md)
- [Deploy to Production](../advanced/deployment.md)
- [Configure Backups](../cli/backup.md)

## See Also

- [Import Command Reference](../cli/import.md)
- [Export Command Reference](../cli/export.md)
- [From R2 Migration](./from-r2.md)
- [SQL Migration](./sql.md)
