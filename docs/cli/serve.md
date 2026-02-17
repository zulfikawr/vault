# vault serve

Start the Vault HTTP server.

## Usage

```bash
vault serve [options]
```

## Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--port` | int | 8090 | Server port |
| `--dir` | string | `./vault_data` | Data directory |
| `--db-path` | string | `{dir}/vault.db` | Database path |
| `--log-level` | string | INFO | Log level (DEBUG/INFO/WARN/ERROR) |
| `--log-format` | string | text | Log format (text/json) |
| `--tls-cert` | string | - | TLS certificate path |
| `--tls-key` | string | - | TLS key path |
| `--jwt-secret` | string | - | JWT secret |
| `--cors-origins` | string | * | CORS origins |
| `--rate-limit` | int | 300 | Rate limit per minute |
| `--max-upload-size` | string | 10MB | Max upload size |
| `--config` | string | config.json | Config file path |

## Examples

### Basic Start

```bash
vault serve
```

### Custom Port

```bash
vault serve --port 3000
```

### Production Mode

```bash
vault serve \
  --port 8090 \
  --log-level WARN \
  --log-format json \
  --rate-limit 100
```

### With TLS

```bash
vault serve \
  --tls-cert /path/to/cert.pem \
  --tls-key /path/to/key.pem
```

### Custom Data Directory

```bash
vault serve --dir /var/lib/vault
```

## Output

```
✓ Vault server started successfully
  Web UI:  http://localhost:8090/
  API:     http://localhost:8090/api

✓ Database initialized at vault_data/vault.db
✓ Admin dashboard available at http://localhost:8090/_/
```

## Access Points

| Service | URL | Description |
|---------|-----|-------------|
| Web UI | http://localhost:8090/ | Main website |
| Admin Dashboard | http://localhost:8090/_/ | Admin interface |
| API | http://localhost:8090/api | REST API |
| Health Check | http://localhost:8090/api/health | Health status |

## Environment Variables

All options can be set via environment:

```bash
VAULT_PORT=8090 \
VAULT_DATA_DIR=./vault_data \
VAULT_LOG_LEVEL=INFO \
VAULT_JWT_SECRET=your-secret \
vault serve
```

## Graceful Shutdown

Press `Ctrl+C` to stop the server. Vault will:
1. Stop accepting new connections
2. Wait for active requests to complete (max 30 seconds)
3. Close database connections
4. Exit cleanly

## Production Deployment

### Systemd Service

```ini
# /etc/systemd/system/vault.service
[Unit]
Description=Vault Backend Server
After=network.target

[Service]
Type=simple
User=vault
WorkingDirectory=/opt/vault
ExecStart=/opt/vault/vault serve --port 8090
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### Docker

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o vault ./cmd/vault

FROM alpine:latest
WORKDIR /vault
COPY --from=builder /app/vault .
EXPOSE 8090
CMD ["./vault", "serve"]
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8090
lsof -i :8090

# Use different port
vault serve --port 8091
```

### Database Locked

```bash
# Check for other processes
lsof vault_data/vault.db

# Remove WAL files if crashed
rm vault_data/vault.db-wal vault_data/vault.db-shm
```

### TLS Issues

```bash
# Generate self-signed cert
openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem -out cert.pem -days 365 -nodes

# Start with TLS
vault serve --tls-cert cert.pem --tls-key key.pem
```

## Performance Tuning

### Increase Rate Limit

```bash
vault serve --rate-limit 1000
```

### Adjust Upload Size

```bash
vault serve --max-upload-size 100MB
```

### Enable Debug Logging

```bash
vault serve --log-level DEBUG --log-format text
```

## Monitoring

### Health Endpoints

```bash
# Basic health check
curl http://localhost:8090/api/health

# Detailed health
curl http://localhost:8090/api/health/collections
```

### Logs

```bash
# View logs in real-time
tail -f vault_data/vault.log

# JSON logs (for log aggregation)
vault serve --log-format json
```

## Next Steps

- [Admin Dashboard](../api/auth.md) - Access the admin UI
- [API Reference](../api/rest.md) - Use the REST API
- [Deployment](../advanced/deployment.md) - Production deployment

## See Also

- [`vault init`](./init.md) - Initialize Vault
- [`vault admin`](./admin.md) - Manage users
- [Configuration](../advanced/config.md) - Advanced config
