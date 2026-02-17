# Migrate from R2 Storage to Vault

This guide covers migrating your Cloudflare R2 bucket to Vault's local storage system.

## Overview

| Source | Destination | Method |
|--------|-------------|--------|
| R2 Bucket | Vault Storage | Direct file copy |
| R2 Objects | Local Files | Preserved paths |

## Prerequisites

- R2 bucket exported (see [From D1 Migration](./from-d1.md#step-2-export-r2-bucket))
- Vault initialized
- Admin credentials

## Step 1: Export R2 Bucket

### Using Wrangler Script

Create `download-r2.js`:

```javascript
const { execSync } = require("child_process");
const { mkdirSync, existsSync, writeFileSync } = require("fs");
const { join, dirname } = require("path");

const BUCKET = "zulfikar-storage";
const OUTPUT_DIR = "./r2-export";
const ACCOUNT_ID = "your-account-id";
const API_TOKEN = "your-api-token";

async function listObjects() {
  const objects = [];
  let cursor = undefined;
  
  do {
    let url = `https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/r2/buckets/${BUCKET}/objects?limit=1000`;
    if (cursor) url += `&cursor=${cursor}`;
    
    const response = await fetch(url, {
      headers: {
        "Authorization": `Bearer ${API_TOKEN}`,
        "Content-Type": "application/json"
      }
    });
    
    const data = await response.json();
    objects.push(...data.result.objects || data.result);
    cursor = data.result.cursor;
  } while (cursor);
  
  return objects;
}

async function downloadObject(key) {
  const url = `https://${ACCOUNT_ID}.r2.cloudflarestorage.com/${BUCKET}/${encodeURIComponent(key)}`;
  
  const response = await fetch(url, {
    headers: {
      "Authorization": `Bearer ${API_TOKEN}`
    }
  });
  
  return response.body;
}

async function main() {
  console.log("📦 Starting R2 download...");
  
  const objects = await listObjects();
  console.log(`Found ${objects.length} objects\n`);
  
  for (const obj of objects) {
    const key = obj.key;
    const outputPath = join(OUTPUT_DIR, key);
    
    const dir = dirname(outputPath);
    if (!existsSync(dir)) {
      mkdirSync(dir, { recursive: true });
    }
    
    console.log(`⬇️  ${key}`);
    
    const body = await downloadObject(key);
    const fileStream = require("fs").createWriteStream(outputPath);
    
    const reader = body.getReader();
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      fileStream.write(value);
    }
    
    fileStream.end();
    console.log(`   ✓ Saved`);
  }
  
  console.log(`\n✅ Complete!`);
}

main().catch(console.error);
```

Run:
```bash
node download-r2.js
```

### Using Wrangler CLI Directly

For small buckets:

```bash
# Create export directory
mkdir -p r2-export

# List objects
npx wrangler r2 bucket info your-bucket

# Download individual objects
npx wrangler r2 object get "your-bucket/path/to/file" \
  --file "./r2-export/path/to/file" --remote
```

### Using rclone (Recommended for Large Buckets)

```bash
# Install rclone
curl https://rclone.org/install.sh | sudo bash

# Configure R2
rclone config

# Select:
# 1. New remote
# Name: r2-remote
# Storage: S3 Compliant Storage
# Provider: Cloudflare R2
# access_key_id: your-access-key
# secret_access_key: your-secret-key
# region: auto
# endpoint: your-account-id.r2.cloudflarestorage.com

# Download entire bucket
rclone copy r2-remote:your-bucket ./r2-export/ --progress
```

## Step 2: Verify Export

```bash
# Count files
find r2-export -type f | wc -l

# Check total size
du -sh r2-export/

# Preview structure
tree -L 2 r2-export/
```

**Expected output:**
```
r2-export/
├── certificates/
│   ├── cs50x-2024-computer-science/
│   └── ...
├── employments/
│   └── ...
└── incremental-cache/
    └── ...

20 files
9.9 MB total
```

## Step 3: Copy to Vault Storage

```bash
# Copy all files
cp -r r2-export/* vault_data/storage/

# Or copy specific folders
cp -r r2-export/certificates vault_data/storage/
cp -r r2-export/employments vault_data/storage/
```

## Step 4: Verify Import

### Using CLI

```bash
# List root storage
vault storage list \
  --path "/" \
  --email "admin@example.com" \
  --password "yourpassword"

# List specific folder
vault storage list \
  --path "/certificates" \
  --recursive \
  --email "admin@example.com" \
  --password "yourpassword"
```

**Expected output:**
```
Storage: /
┌─────────────────────────────────────┬──────────┬──────────────┐
│ Path                                │ Type     │ Size         │
├─────────────────────────────────────┼──────────┼──────────────┤
│ certificates/                       │ folder   │ -            │
│ employments/                        │ folder   │ -            │
│ incremental-cache/                  │ folder   │ -            │
└─────────────────────────────────────┴──────────┴──────────────┘
```

### Check File System

```bash
# List storage directory
ls -la vault_data/storage/

# Check file count
find vault_data/storage -type f | wc -l

# Verify sizes match
du -sh vault_data/storage/
```

## Step 5: Test File Access

### Via API

```bash
# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8090/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' \
  | jq -r '.token')

# List files
curl -s "http://localhost:8090/api/storage?path=/" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Download file
curl -s "http://localhost:8090/api/storage/certificates/cs50x-2024-computer-science/image.png" \
  -H "Authorization: Bearer $TOKEN" \
  -o downloaded-image.png

# Verify download
ls -lh downloaded-image.png
```

### Via Admin Dashboard

1. Open http://localhost:8090/_/
2. Navigate to "Storage"
3. Browse folders
4. Click files to download

## File Path Mapping

### R2 → Vault Storage

| R2 Path | Vault Storage Path | URL Path |
|---------|-------------------|----------|
| `certificates/image.png` | `storage/certificates/image.png` | `/api/storage/certificates/image.png` |
| `employments/logo.png` | `storage/employments/logo.png` | `/api/storage/employments/logo.png` |
| `a/b/c/file.txt` | `storage/a/b/c/file.txt` | `/api/storage/a/b/c/file.txt` |

### Update Your Application

**Before (R2):**
```javascript
const imageUrl = `https://r2.yourdomain.com/certificates/${id}/image.png`;
```

**After (Vault):**
```javascript
const imageUrl = `http://localhost:8090/api/storage/certificates/${id}/image.png`;

// Or with auth token
const imageUrl = `http://localhost:8090/api/storage/certificates/${id}/image.png`;
const headers = {
  "Authorization": `Bearer ${token}`
};
```

## Handling Large Files

### Multipart Upload Support

Vault supports multipart uploads for files >5MB:

```bash
# Upload large file
curl -X POST http://localhost:8090/api/storage \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@large-video.mp4" \
  -F "path=/videos/"
```

### Storage Limits

Default upload size: 10MB

Increase in `config.json`:
```json
{
  "max_file_upload_size": 104857600  // 100MB
}
```

## Storage Structure

### Recommended Organization

```
vault_data/storage/
├── users/              # User uploads
│   └── {user-id}/
├── collections/        # Collection files
│   └── {collection-name}/
├── cache/             # Temporary cache
└── backups/           # File backups
```

### Access Control

Set authorization rules for file access:

```bash
# Via API
curl -X PATCH http://localhost:8090/api/collections/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "create_rule": "@request.auth.id != \"\"",
    "update_rule": "@request.auth.id = record.author_id"
  }'
```

## Troubleshooting

### "File not found"

Check the path:
```bash
# List to verify file exists
vault storage list --path "/certificates" --email "admin@example.com" --password "yourpassword"
```

### Permission Denied

Check file permissions:
```bash
ls -la vault_data/storage/
chmod -R 755 vault_data/storage/
```

### Missing Files

Compare counts:
```bash
# R2 count
rclone lsf r2-remote:your-bucket | wc -l

# Vault count
find vault_data/storage -type f | wc -l
```

### Path Traversal Protection

Vault blocks `..` in paths for security. Use absolute paths:
```bash
# ✅ Correct
vault storage create --path "/uploads/file.png" --file "./file.png"

# ❌ Incorrect
vault storage create --path "../file.png" --file "./file.png"
```

## Post-Migration

### Cleanup

```bash
# Remove export directory
rm -rf r2-export/

# Or archive for backup
tar -czf r2-backup-$(date +%Y%m%d).tar.gz r2-export/
```

### Configure Backups

```bash
# Create backup including storage
vault backup create --output "full-backup.zip"

# Verify backup includes storage
unzip -l full-backup.zip | grep storage
```

### Set Up Monitoring

Monitor storage usage:
```bash
# Check storage size
du -sh vault_data/storage/

# Check file count
find vault_data/storage -type f | wc -l
```

## Performance Tips

### Optimize File Access

1. **Use CDN**: Put Vault behind a CDN for static assets
2. **Cache Headers**: Configure appropriate cache headers
3. **Compression**: Enable gzip for text files

### Storage Best Practices

1. **Organize by Collection**: Keep files with their collections
2. **Use Meaningful Paths**: `users/{id}/avatar.png` not `files/123.png`
3. **Regular Cleanup**: Remove orphaned files periodically

## Security Considerations

### File Upload Validation

- Validate file types
- Limit file sizes
- Scan for malware
- Sanitize filenames

### Access Control

- Require authentication for file access
- Use authorization rules
- Implement rate limiting

## Next Steps

- [Configure File Upload](../api/files.md)
- [Set Up Authorization](../concepts/rules.md)
- [Deploy to Production](../advanced/deployment.md)

## See Also

- [From D1 Migration](./from-d1.md)
- [Storage CLI Reference](../cli/storage.md)
- [File API](../api/files.md)
