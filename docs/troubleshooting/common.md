# Common Issues

## Server Won't Start

### Port Already in Use
```bash
# Find process
lsof -i :8090

# Use different port
vault serve --port 8091
```

### Database Locked
```bash
# Remove WAL files
rm vault_data/vault.db-wal vault_data/vault.db-shm
```

## Authentication Issues

### "Invalid credentials"
- Check email/password
- Reset password: `vault admin reset-password`

### "Token expired"
- Login again to get new token
- Increase expiry: `VAULT_JWT_EXPIRY=168`

## Import/Export Issues

### "Collection not found"
Create collection first or use `--collection` flag.

### "Table already exists"
Run `vault migrate sync` to register.

## Performance Issues

### Slow queries
- Add indexes on filtered fields
- Use pagination

### High memory usage
- Reduce `perPage` limit
- Enable log rotation

See Also: [FAQ](./faq.md), [Error Codes](./errors.md)
