# CLI Overview

The Vault CLI provides a comprehensive set of commands for managing your Vault instance.

## Usage Pattern

```bash
vault <command> [subcommand] [options]
```

## Command Categories

### Project Management
- [`vault init`](./init.md) - Initialize new Vault project
- [`vault serve`](./serve.md) - Start the HTTP server
- [`vault version`](#version) - Display version information

### User Management
- [`vault admin`](./admin.md) - Manage admin users

### Schema Management
- [`vault collection`](./collection.md) - Manage collections
- [`vault migrate`](./migrate.md) - Database migrations

### Data Management
- [`vault export`](./export.md) - Export data
- [`vault import`](./import.md) - Import data
- [`vault backup`](./backup.md) - Backup and restore

### File Management
- [`vault storage`](./storage.md) - Manage file storage

## Global Options

| Option | Description |
|--------|-------------|
| `--config FILE` | Path to config file (default: `config.json`) |
| `-h, --help` | Show help for command |
| `--env ENV` | Environment name |
| `--env-file FILE` | Path to .env file |

## Environment Variables

All configuration can be overridden using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `VAULT_PORT` | Server port | 8090 |
| `VAULT_DATA_DIR` | Data directory | `./vault_data` |
| `VAULT_DB_PATH` | Database path | `{data_dir}/vault.db` |
| `VAULT_LOG_LEVEL` | Log level | INFO |
| `VAULT_LOG_FORMAT` | Log format | text |
| `VAULT_JWT_SECRET` | JWT secret | (auto-generated) |
| `VAULT_JWT_EXPIRY` | JWT expiry (hours) | 72 |
| `VAULT_CORS_ORIGINS` | CORS origins | * |
| `VAULT_RATE_LIMIT_PER_MIN` | Rate limit | 300 |
| `VAULT_MAX_FILE_UPLOAD_SIZE` | Max upload size | 10MB |

## Examples

### Get Help

```bash
# General help
vault -h

# Command-specific help
vault init -h
vault collection -h
```

### Common Workflows

```bash
# Initialize new project
vault init --email "admin@example.com" --username "admin" --password "secret"

# Start server
vault serve --port 8090

# Create collection
vault collection create --name "posts" --fields "title:text,body:text" \
  --email "admin@example.com" --password "secret"

# List collections
vault collection list --email "admin@example.com" --password "secret"

# Create backup
vault backup create --output "backup.zip"

# Export data
vault export json --output "backup.json"
```

## Configuration File

Vault uses `config.json` for configuration:

```json
{
  "port": 8090,
  "db_path": "./vault_data/vault.db",
  "data_dir": "./vault_data",
  "log_level": "INFO",
  "log_format": "text",
  "jwt_secret": "your-secret-key",
  "jwt_expiry": 72,
  "max_file_upload_size": 10485760,
  "cors_origins": "*",
  "rate_limit_per_min": 300
}
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Configuration error |
| 4 | Database error |
| 5 | Authentication error |

## Output Formats

Most commands use human-readable output by default. Some commands support:

- `--json` - JSON output (for scripting)
- `--quiet` - Minimal output

## Next Steps

- [vault init](./init.md) - Start a new project
- [vault serve](./serve.md) - Run the server
- [API Reference](../api/rest.md) - Use the REST API
