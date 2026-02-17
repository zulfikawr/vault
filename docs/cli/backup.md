# vault backup

Backup and restore Vault data.

## Usage

```bash
vault backup <subcommand> [options]
```

## Subcommands

### create

Create a backup of database and storage.

```bash
vault backup create [--output FILE]
```

**Options:**
- `--output`: Output file path (default: `vault_backup_TIMESTAMP.zip`)

**Example:**
```bash
# Create backup with timestamp
vault backup create

# Create backup with custom name
vault backup create --output "./my-backup.zip"
```

**Backup Contents:**
- `vault.db` - SQLite database
- `config.json` - Configuration file
- `storage/` - All uploaded files

**Output:**
```
✓ Backup created successfully
  File: vault_backup_20260217_120000.zip
  Size: 15.5 MB
```

### list

List all backups.

```bash
vault backup list
```

**Output:**
```
Total backups: 3

Filename                                 Size            Modified            
-----------------------------------------------------------
vault_backup_20260217_120000.zip         15.50 MB        2026-02-17 12:00:00
vault_backup_20260216_100000.zip         14.20 MB        2026-02-16 10:00:00
vault_backup_20260215_080000.zip         13.80 MB        2026-02-15 08:00:00
```

### restore

Restore from a backup.

```bash
vault backup restore --input FILE [--force]
```

**Options:**
- `--input` (required): Backup file path
- `--force`: Skip confirmation prompt

**Example:**
```bash
# Restore with confirmation
vault backup restore --input "./vault_backup_20260217_120000.zip"

# Restore without confirmation
vault backup restore --input "./backup.zip" --force
```

**Warning:** This will overwrite your current database and storage!

**Output:**
```
This will overwrite your current database and storage. Continue? (yes/no): yes
✓ Backup restored successfully
  From: vault_backup_20260217_120000.zip
```

## Backup Format

Backups are ZIP files containing:

```
vault_backup_TIMESTAMP.zip
├── vault.db           # SQLite database
├── config.json        # Configuration
└── storage/           # File storage
    ├── folder1/
    └── folder2/
```

## Automated Backups

### Daily Backup Script

```bash
#!/bin/bash
# /usr/local/bin/vault-backup-daily.sh

BACKUP_DIR="/var/backups/vault"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=7

mkdir -p $BACKUP_DIR

# Create backup
cd /opt/vault
vault backup create --output "$BACKUP_DIR/vault_$TIMESTAMP.zip"

# Clean old backups
find $BACKUP_DIR -name "vault_*.zip" -mtime +$KEEP_DAYS -delete
```

### Cron Job

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /usr/local/bin/vault-backup-daily.sh
```

## Restore Procedure

1. **Stop Vault server:**
   ```bash
   pkill vault
   ```

2. **Restore backup:**
   ```bash
   vault backup restore --input "./backup.zip" --force
   ```

3. **Start Vault server:**
   ```bash
   vault serve
   ```

4. **Verify:**
   ```bash
   vault collection list --email "admin@example.com" --password "secret"
   ```

## Best Practices

1. **Regular backups**: Daily or before major changes
2. **Off-site storage**: Copy backups to different location
3. **Test restores**: Periodically test backup restoration
4. **Retention policy**: Keep 7-30 days of backups
5. **Secure backups**: Encrypt sensitive backups

## Troubleshooting

### "Backup file not found"

Check the file exists:
```bash
ls -la vault_backup_*.zip
```

### "Restore failed - file corrupted"

Verify ZIP integrity:
```bash
unzip -t vault_backup_TIMESTAMP.zip
```

### "Not enough disk space"

Check available space:
```bash
df -h
```

## See Also

- [`vault export`](./export.md) - Export data
- [`vault import`](./import.md) - Import data
- [Deployment](../advanced/deployment.md) - Production setup
