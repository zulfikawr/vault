# Import Command

Import data into Vault from external sources including Cloudflare D1, JSON, and SQL files.

## Usage

```bash
vault import <format> [options] <file>
```

## Formats

### Import from Cloudflare D1

Import data exported from Cloudflare D1 databases.

#### Export from D1 First

```bash
# Export your D1 database
npx wrangler d1 export <DATABASE_NAME> --output d1-export.sql --remote
```

#### Import to Vault

```bash
vault import d1 d1-export.sql
```

#### Dry Run (Validate Without Importing)

```bash
vault import d1 d1-export.sql --dry-run
```

**What happens during D1 import:**
1. Creates tables from D1 SQL dump
2. Imports all records (handles ~3,500+ records)
3. Automatically registers tables as Vault collections
4. Maps SQL types to Vault field types:
   - `TEXT` → `text`
   - `INTEGER` → `boolean` (if contains "BOOL")
   - `REAL/DOUBLE/FLOAT` → `number`

### Import from JSON

Import from Vault JSON export or simple record arrays.

#### Import Vault Export

```bash
# Import full export (includes schema)
vault import json backup.json
```

#### Import Simple Records Array

```bash
# Import to existing collection
vault import json users.json --collection users
```

**Simple JSON format:**
```json
[
  {
    "id": "usr_123",
    "email": "admin@example.com",
    "created": "2026-02-17T12:00:00Z"
  }
]
```

#### Dry Run

```bash
vault import json data.json --collection users --dry-run
```

### Import from SQL

Import from generic SQL files.

```bash
vault import sql schema.sql
```

**Supported SQL statements:**
- `CREATE TABLE`
- `INSERT INTO`

## Options

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `--dry-run` | All | Validate without importing | `false` |
| `--collection NAME` | JSON | Target collection (for simple JSON) | Auto-detect |

## Examples

### Migrate from Cloudflare D1

```bash
# Step 1: Export D1
npx wrangler d1 export homepage-db --output d1.sql --remote

# Step 2: Initialize Vault
vault init --email "admin@example.com" --username "admin" --password "secret"

# Step 3: Import D1 data
vault import d1 d1.sql

# Step 4: Verify import
vault collection list --email "admin@example.com" --password "secret"
```

### Import User Data

```bash
# Export users from another system
# ... (export process)

# Import to Vault users collection
vault import json users-export.json --collection users
```

### Test Import Before Production

```bash
# Test in development first
vault import d1 production-dump.sql --dry-run

# If successful, import for real
vault import d1 production-dump.sql
```

## Migration Workflow: D1 → Vault

### Complete Example

```bash
# 1. Export from Cloudflare
npx wrangler d1 export mydb --output mydb.sql --remote

# 2. Initialize fresh Vault
rm -rf vault_data config.json
vault init --email "admin@example.com" --username "admin" --password "secret"

# 3. Import with dry-run first
vault import d1 mydb.sql --dry-run

# 4. Import for real
vault import d1 mydb.sql

# 5. Verify collections
vault collection list --email "admin@example.com" --password "secret"

# 6. Check record counts via API
curl -s http://localhost:8090/api/collections/posts/records?page=1&perPage=1 \
  | jq '.totalRecords'
```

### Expected Output

```
📥 Importing D1 dump: mydb.sql
Found 3550 SQL statements

✓ Created table: users
✓ Created table: analytics_events
  ... 500 records imported
  ... 1000 records imported
  ... 1500 records imported
  ... 2000 records imported
  ... 2500 records imported
  ... 3000 records imported
✓ Created table: posts
✓ Created table: projects
✓ Created table: certificates

✅ D1 Import complete!
   Tables: 17
   Records: 3528
   Errors: 1

📋 Registering collections in Vault...
   Found 17 user tables: [users posts projects ...]
   ✓ Registered: users (4 fields)
   ✓ Registered: posts (9 fields)
   ✓ Registered: projects (13 fields)
   Registered 17 collections
```

## Field Type Mapping

### D1/SQLite → Vault

| SQL Type | Vault Type | Notes |
|----------|-----------|-------|
| `TEXT` | `text` | Default type |
| `INTEGER` | `boolean` | If column name contains "bool" |
| `INTEGER` | `number` | Otherwise |
| `REAL` | `number` | Floating point numbers |
| `DOUBLE` | `number` | Double precision |
| `FLOAT` | `number` | Single precision |
| `BLOB` | `text` | Stored as base64 |
| `DATE/TIME` | `text` | ISO 8601 format |

### Standard Fields

All imported collections include Vault's standard fields:
- `id` (TEXT PRIMARY KEY)
- `created` (TEXT) - Record creation timestamp
- `updated` (TEXT) - Last update timestamp

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `IMPORT_FORMAT_REQUIRED` | 400 | No import format specified |
| `INVALID_IMPORT_FORMAT` | 400 | Unknown import format |
| `INVALID_FLAGS` | 400 | Failed to parse flags |
| `SQL_FILE_REQUIRED` | 400 | SQL file path required |
| `JSON_FILE_REQUIRED` | 400 | JSON file path required |
| `D1_FILE_REQUIRED` | 400 | D1 file path required |
| `FILE_NOT_FOUND` | 404 | Input file not found |
| `FILE_READ_FAILED` | 500 | Failed to read file |
| `JSON_PARSE_FAILED` | 400 | Failed to parse JSON |
| `COLLECTION_REQUIRED` | 400 | Collection flag required |
| `COLLECTION_NOT_FOUND` | 404 | Target collection doesn't exist |
| `DB_CONNECTION_FAILED` | 500 | Failed to connect to database |
| `SCHEMA_LOAD_FAILED` | 500 | Failed to load schema |

## Troubleshooting

### "Table already exists"

Tables are created but not registered. Run:

```bash
vault migrate sync
```

### "No collections found" after import

The collections are in the database but not showing in `migrate status`. Use:

```bash
vault collection list --email "admin@example.com" --password "yourpassword"
```

### Import fails with "database is locked"

Make sure no other process is using the database:

```bash
# Stop any running Vault server
pkill vault

# Try import again
vault import d1 dump.sql
```

### Missing records after import

Check for errors in the output. Some records may fail due to:
- Duplicate IDs
- Constraint violations
- Invalid data types

Re-import with dry-run to identify issues:

```bash
vault import d1 dump.sql --dry-run
```

### Path Issues

Make sure you're in the directory with `config.json`:

```bash
# Check current directory
pwd
ls config.json

# Or specify data directory
VAULT_DATA_DIR=./vault_data vault import d1 dump.sql
```

## Best Practices

1. **Always Dry-Run First**: Test imports before running them
2. **Backup Before Import**: Export current data before importing
3. **Test in Development**: Never import directly to production
4. **Verify After Import**: Check record counts and data integrity
5. **Clean Up Duplicates**: Handle duplicate IDs before importing

## Post-Import Checklist

- [ ] Verify collection count: `vault collection list`
- [ ] Check record counts via API
- [ ] Test authentication (if users imported)
- [ ] Verify file storage (if applicable)
- [ ] Test API endpoints
- [ ] Review authorization rules
- [ ] Update application configuration

## See Also

- [Export Command](./export.md)
- [From D1 Migration](../migration/from-d1.md)
- [SQL Migration](../migration/sql.md)
- [Backup Command](./backup.md)
