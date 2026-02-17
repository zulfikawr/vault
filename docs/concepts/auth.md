# Authentication

Vault uses JWT-based authentication with bcrypt password hashing.

## Login

```bash
curl -X POST http://localhost:8090/api/collections/users/auth-with-password \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "secret"}'
```

**Response:**
```json
{
  "token": "eyJhbGc...",
  "record": {"id": "usr_123", "email": "admin@example.com"}
}
```

## Refresh Token

```bash
curl -X POST http://localhost:8090/api/collections/users/auth-refresh \
  -H "Authorization: Bearer TOKEN"
```

## Password Reset

1. Request reset:
```bash
curl -X POST http://localhost:8090/api/collections/users/request-password-reset \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com"}'
```

2. Confirm reset:
```bash
curl -X POST http://localhost:8090/api/collections/users/confirm-password-reset \
  -H "Content-Type: application/json" \
  -d '{"token": "TOKEN", "password": "newpassword"}'
```

## Token Structure

```json
{
  "record_id": "usr_123",
  "collection": "users",
  "request_id": "req_abc",
  "exp": 1234567890
}
```

## Security

- Passwords hashed with bcrypt
- Tokens expire after 72 hours (configurable)
- Refresh tokens valid for 30 days

See Also: [API Auth](../api/auth.md)
