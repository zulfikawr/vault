# vault migrate

Database migration operations.

## Usage

```bash
vault migrate <subcommand> [options]
```

## Subcommands

### sync

Synchronize database schema with collections.

```bash
vault migrate sync [--collection NAME] [--verbose]
```

**Options:**
- `--collection`: Sync specific collection only
- `--verbose`: Verbose output

**Example:**
```bash
# Sync all collections
vault migrate sync

# Sync specific collection
vault migrate sync --collection "posts" --verbose
```

**Output:**
```
Syncing 4 collection(s)...

✓ _audit_logs
✓ _collections
✓ users
✓ _refresh_tokens

Migration Summary:
  Total: 4
  Success: 4
  Failed: 0

✓ All collections synced successfully
```

### status

Show current database and collection status.

```bash
vault migrate status
```

**Output:**
```
Database Status:
  Database: ./vault_data/vault.db
  Collections: 4

Collection                     Type            Fields    
-----------------------------------------------------------
_collections                   system          8         
_audit_logs                    system          5         
_refresh_tokens                system          3         
users                          auth            4         
```

## What Sync Does

1. **Checks table existence**: Verifies each collection has a corresponding table
2. **Creates missing tables**: Creates tables for new collections
3. **Adds new columns**: Adds fields that don't exist in the table
4. **Creates indexes**: Creates indexes defined on collections
5. **Preserves data**: Existing data is not modified

## When to Run Sync

- After creating collections via API
- After modifying collection schema
- After importing data
- When schema is out of sync

## System Collections

Vault automatically manages these system collections:

| Collection | Type | Description |
|------------|------|-------------|
| `_collections` | system | Collection definitions |
| `_audit_logs` | system | Audit log entries |
| `_refresh_tokens` | system | JWT refresh tokens |
| `users` | auth | User accounts |

## Field Type Mapping

| Vault Type | SQLite Type |
|------------|-------------|
| `text` | TEXT |
| `number` | REAL |
| `bool` | INTEGER |
| `date` | TEXT |
| `json` | TEXT |
| `relation` | TEXT |
| `file` | TEXT |

## Troubleshooting

### "Collection not found"

Check collection exists:
```bash
vault collection list --email "admin@example.com" --password "secret"
```

### "Table already exists"

This is normal for existing collections. Sync will update schema if needed.

### "Failed to sync"

Run with verbose to see details:
```bash
vault migrate sync --verbose
```

## See Also

- [`vault collection`](./collection.md) - Manage collections
- [`vault import`](./import.md) - Import data
- [Schema](../concepts/collections.md) - Collection schema
