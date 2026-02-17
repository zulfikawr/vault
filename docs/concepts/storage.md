# Storage

File storage system with local filesystem backend.

## Storage Structure

```
vault_data/storage/
├── collection-name/
│   └── record-id/
│       └── file.png
```

## Upload Files

```bash
curl -X POST http://localhost:8090/api/files \
  -H "Authorization: Bearer TOKEN" \
  -F "file=@image.png" \
  -F "path=/uploads/"
```

## Download Files

```bash
curl http://localhost:8090/api/files/collection/record-id/file.png \
  -H "Authorization: Bearer TOKEN"
```

## Limits

- Default max upload: 10MB
- Configurable via `max_file_upload_size`

See Also: [Storage CLI](../cli/storage.md)
