# FAQ

## General

### What is Vault?
A self-contained backend framework with dynamic schema, auth, and real-time features.

### Is it production-ready?
Yes, but review security settings before deploying.

### What database does it use?
SQLite (embedded, no server needed).

## Development

### How do I reset admin password?
```bash
vault admin reset-password --email "admin@example.com" --password "newpass"
```

### Can I use Vault with existing database?
Use `vault import` to migrate data.

### How do I enable HTTPS?
```bash
vault serve --tls-cert cert.pem --tls-key key.pem
```

## Deployment

### How do I deploy to production?
See [Deployment Guide](../advanced/deployment.md).

### Can I scale horizontally?
Vault is designed for single-instance deployment. For scaling, use load balancer with shared storage.

### How do I backup data?
```bash
vault backup create --output backup.zip
```

See Also: [Common Issues](./common.md), [Error Codes](./errors.md)
