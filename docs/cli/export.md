# Export Command

Export collections and data from Vault in various formats.

## Usage

```bash
vault export <format> [options]
```

## Formats

### JSON Export

Export collections and records as JSON.

#### Export All Collections

```bash
vault export json --output ./backup.json
```

#### Export Specific Collection

```bash
vault export json --collection users --output ./users.json
```

#### Pretty Print JSON

```bash
vault export json --collection users --pretty
```

#### Output Format

```json
{
  "exported_at": "2026-02-17T12:00:00Z",
  "vault_version": "0.7.0",
  "collections": {
    "users": {
      "schema": {
        "id": "col_users",
        "name": "users",
        "type": "auth",
        "fields": [
          {"name": "email", "type": "text"},
          {"name": "password", "type": "text"}
        ]
      },
      "records": [
        {
          "id": "usr_123",
          "created": "2026-02-17T12:00:00Z",
          "updated": "2026-02-17T12:00:00Z",
          "email": "admin@example.com"
        }
      ],
      "count": 1
    }
  }
}
```

### SQL Export

Export schema and data as SQL statements.

```bash
vault export sql --output ./schema.sql
```

#### Output Format

```sql
-- Vault SQL Export
-- Generated: 2026-02-17T12:00:00Z
-- Vault Version: 0.7.0

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  created TEXT,
  updated TEXT,
  email TEXT,
  password TEXT
);

INSERT INTO users (id, created, updated, email, password) VALUES (?, ?, ?, ?, ?);
```

## Options

| Option | Format | Description | Default |
|--------|--------|-------------|---------|
| `--output FILE` | All | Output file path | `vault_export_TIMESTAMP.ext` |
| `--collection NAME` | JSON | Export specific collection only | All collections |
| `--pretty` | JSON | Pretty print JSON output | `false` |

## Examples

### Backup Before Migration

```bash
# Export everything
vault export json --output ./pre-migration-backup.json
vault export sql --output ./pre-migration-backup.sql
```

### Export Single Collection for Sharing

```bash
vault export json --collection posts --pretty --output ./posts-export.json
```

### Automated Backup Script

```bash
#!/bin/bash
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

vault export json --output "$BACKUP_DIR/vault_$TIMESTAMP.json"
vault export sql --output "$BACKUP_DIR/vault_$TIMESTAMP.sql"

# Keep only last 7 backups
find $BACKUP_DIR -name "vault_*.json" -mtime +7 -delete
find $BACKUP_DIR -name "vault_*.sql" -mtime +7 -delete
```

## Import Exported Data

See [Import Command](./import.md) for importing exported data.

```bash
# Import JSON export
vault import json ./backup.json

# Import specific collection from JSON
vault import json ./users.json --collection users
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `EXPORT_FORMAT_REQUIRED` | 400 | No export format specified |
| `INVALID_EXPORT_FORMAT` | 400 | Unknown export format |
| `DB_CONNECTION_FAILED` | 500 | Failed to connect to database |
| `SCHEMA_LOAD_FAILED` | 500 | Failed to load schema |
| `COLLECTION_NOT_FOUND` | 404 | Specified collection doesn't exist |
| `JSON_MARSHAL_FAILED` | 500 | Failed to marshal JSON |
| `FILE_WRITE_FAILED` | 500 | Failed to write output file |

## Troubleshooting

### "No collections found to export"

Make sure you have created collections:

```bash
vault collection list --email "admin@example.com" --password "yourpassword"
```

### "Collection not found"

Check the collection name:

```bash
vault collection list --email "admin@example.com" --password "yourpassword"
```

### Permission Denied

Make sure you have write permissions in the output directory:

```bash
ls -la ./output-directory/
```

## Best Practices

1. **Regular Backups**: Export data regularly before making schema changes
2. **Version Control**: Commit SQL exports to version control for schema tracking
3. **Secure Storage**: Store JSON exports securely as they contain all data
4. **Test Imports**: Always test imports in a development environment first

## See Also

- [Import Command](./import.md)
- [Backup Command](./backup.md)
- [Migration Guide](../migration/sql.md)
