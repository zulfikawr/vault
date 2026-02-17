# Error Codes

## Authentication Errors

| Code | HTTP | Description |
|------|------|-------------|
| `INVALID_TOKEN` | 401 | Token invalid or expired |
| `INVALID_CREDENTIALS` | 401 | Wrong email/password |
| `MISSING_AUTH` | 401 | No auth header |

## Authorization Errors

| Code | HTTP | Description |
|------|------|-------------|
| `ACCESS_DENIED` | 403 | Rule evaluation failed |
| `FORBIDDEN` | 403 | Insufficient permissions |

## Database Errors

| Code | HTTP | Description |
|------|------|-------------|
| `DB_CONNECTION_FAILED` | 500 | Can't connect to DB |
| `SCHEMA_LOAD_FAILED` | 500 | Can't load schema |
| `RECORD_NOT_FOUND` | 404 | Record doesn't exist |

## File Errors

| Code | HTTP | Description |
|------|------|-------------|
| `FILE_NOT_FOUND` | 404 | File doesn't exist |
| `UPLOAD_FAILED` | 500 | Upload error |
| `FILE_TOO_LARGE` | 400 | Exceeds size limit |

## Import/Export Errors

| Code | HTTP | Description |
|------|------|-------------|
| `IMPORT_FORMAT_REQUIRED` | 400 | No format specified |
| `FILE_READ_FAILED` | 500 | Can't read file |
| `JSON_PARSE_FAILED` | 400 | Invalid JSON |
