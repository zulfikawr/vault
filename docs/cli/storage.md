# vault storage

Manage file storage.

## Usage

```bash
vault storage <subcommand> [options]
```

## Subcommands

### create (upload)

Upload a file to storage.

```bash
vault storage create --path PATH --file FILE --email EMAIL --password PASSWORD
```

**Options:**
- `--path` (required): Storage path (e.g., `/uploads/image.png`)
- `--file` (required): Local file to upload
- `--email` (required): Admin email
- `--password` (required): Admin password

**Example:**
```bash
vault storage create --path "/certificates/cert.png" \
  --file "./cert.png" \
  --email "admin@example.com" --password "secret"
```

### list

List files and folders.

```bash
vault storage list --path PATH [--recursive] --email EMAIL --password PASSWORD
```

**Options:**
- `--path`: Path to list (default: `/`)
- `--recursive`: List recursively
- `--email`: Admin email
- `--password`: Admin password

**Example:**
```bash
# List root
vault storage list --email "admin@example.com" --password "secret"

# List specific folder
vault storage list --path "/certificates" \
  --email "admin@example.com" --password "secret"

# List recursively
vault storage list --path "/uploads" --recursive \
  --email "admin@example.com" --password "secret"
```

### get (download)

Download a file from storage.

```bash
vault storage get --path PATH --output FILE --email EMAIL --password PASSWORD [--force]
```

**Options:**
- `--path` (required): File path in storage
- `--output` (required): Output file path
- `--email` (required): Admin email
- `--password` (required): Admin password
- `--force`: Overwrite existing file

**Example:**
```bash
vault storage get --path "/certificates/cert.png" \
  --output "./downloaded-cert.png" \
  --email "admin@example.com" --password "secret"
```

### delete

Delete a file or folder.

```bash
vault storage delete --path PATH [--recursive] [--force] --email EMAIL --password PASSWORD
```

**Options:**
- `--path` (required): Path to delete
- `--recursive`: Delete directory recursively
- `--force`: Skip confirmation
- `--email`: Admin email
- `--password`: Admin password

**Example:**
```bash
# Delete single file
vault storage delete --path "/uploads/temp.png" \
  --email "admin@example.com" --password "secret"

# Delete folder recursively
vault storage delete --path "/uploads/old" --recursive --force \
  --email "admin@example.com" --password "secret"
```

## Storage Structure

```
vault_data/storage/
├── certificates/
│   ├── cert1.png
│   └── cert2.png
├── uploads/
│   └── user123/
│       └── avatar.jpg
└── cache/
    └── temp/
```

## File Limits

| Setting | Default | Config Key |
|---------|---------|------------|
| Max Upload Size | 10MB | `max_file_upload_size` |
| Allowed Paths | All | - |
| Path Traversal | Blocked | Security feature |

## Security

- Path traversal (`..`) is blocked
- Authentication required for all operations
- Files stored with original names
- No file type validation by default

## Best Practices

1. **Organize by collection**: `/collection-name/record-id/file.png`
2. **Use meaningful names**: Avoid generic names like `file1.png`
3. **Clean up old files**: Regularly remove unused files
4. **Set size limits**: Configure `max_file_upload_size` appropriately

## Troubleshooting

### "File already exists"

Use `--force` to overwrite or delete first:
```bash
vault storage delete --path "/file.png" --force --email "admin@example.com" --password "secret"
```

### "Path not found"

Check the path exists:
```bash
vault storage list --path "/" --email "admin@example.com" --password "secret"
```

### "File too large"

Increase upload size in config.json:
```json
{
  "max_file_upload_size": 104857600  // 100MB
}
```

## API Access

Files can be accessed via API:

```bash
# Download file
curl http://localhost:8090/api/storage/path/to/file.png \
  -H "Authorization: Bearer TOKEN"

# List files
curl http://localhost:8090/api/storage?path=/uploads \
  -H "Authorization: Bearer TOKEN"
```

## See Also

- [File Storage](../concepts/storage.md) - Storage concepts
- [File API](../api/files.md) - File endpoints
- [Configuration](../advanced/config.md) - Upload limits
