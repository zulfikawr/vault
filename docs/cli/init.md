# vault init

Initialize a new Vault project with admin user and database.

## Usage

```bash
vault init [options]
```

## Options

| Option | Type | Default | Required | Description |
|--------|------|---------|----------|-------------|
| `--email` | string | - | Yes | Admin email address |
| `--username` | string | - | Yes | Admin username |
| `--password` | string | - | Yes | Admin password |
| `--dir` | string | `./vault_data` | No | Data directory |
| `--skip-admin` | bool | `false` | No | Skip admin creation |
| `--force` | bool | `false` | No | Overwrite existing setup |

## Examples

### Basic Initialization

```bash
vault init --email "admin@example.com" --username "admin" --password "securepassword"
```

### Custom Data Directory

```bash
vault init --email "admin@example.com" --username "admin" \
  --password "securepassword" --dir "./my-vault-data"
```

### Skip Admin Creation

```bash
vault init --skip-admin --dir "./vault-data"

# Create admin later
vault admin create --email "admin@example.com" --username "admin" --password "secret"
```

### Force Re-initialization

```bash
# Warning: This will overwrite existing data!
vault init --email "admin@example.com" --username "admin" \
  --password "newpassword" --force
```

## What It Does

1. **Creates directory structure**:
   ```
   vault_data/
   ├── vault.db          # SQLite database
   └── storage/          # File storage
   ```

2. **Generates configuration**:
   - `config.json` - Main configuration
   - `.env.example` - Environment variables template

3. **Initializes database**:
   - Creates SQLite database
   - Sets up system tables (`_collections`, `_audit_logs`, etc.)

4. **Creates admin user** (unless `--skip-admin`):
   - Hashes password with bcrypt
   - Stores in `users` table

## Output

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

## Configuration Files

### config.json

```json
{
  "port": 8090,
  "db_path": "./vault_data/vault.db",
  "data_dir": "./vault_data",
  "log_level": "INFO",
  "log_format": "text",
  "jwt_secret": "auto-generated-secret",
  "jwt_expiry": 72,
  "max_file_upload_size": 10485760,
  "cors_origins": "*",
  "rate_limit_per_min": 300
}
```

### .env.example

```bash
VAULT_PORT=8090
VAULT_DATA_DIR=./vault_data
VAULT_JWT_SECRET=your-secret-key-here
VAULT_JWT_EXPIRY_HOURS=24
```

## Security Notes

1. **Change the JWT secret** in production
2. **Use strong passwords** (minimum 8 characters)
3. **Store credentials securely** - consider using a secrets manager
4. **Don't commit** `config.json` to version control (contains secrets)

## Troubleshooting

### "Vault already initialized"

Use `--force` to overwrite (warning: deletes existing data):

```bash
vault init --email "admin@example.com" --username "admin" \
  --password "secret" --force
```

### "Invalid email format"

Email must contain `@`:

```bash
# ❌ Wrong
vault init --email "admin" --username "admin" --password "secret"

# ✅ Correct
vault init --email "admin@example.com" --username "admin" --password "secret"
```

### "Password too short"

Password must be at least 8 characters:

```bash
vault init --email "admin@example.com" --username "admin" \
  --password "verysecurepassword123"
```

## Next Steps

1. **Start the server**: [`vault serve`](./serve.md)
2. **Create collections**: [`vault collection create`](./collection.md)
3. **Access dashboard**: http://localhost:8090/_/

## See Also

- [`vault admin`](./admin.md) - Manage admin users
- [`vault serve`](./serve.md) - Start the server
- [Quick Start](../quickstart.md) - Get started guide
