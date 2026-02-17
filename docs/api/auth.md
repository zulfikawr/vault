# Authentication API

## Login

**POST** `/api/collections/users/auth-with-password`

```bash
curl -X POST http://localhost:8090/api/collections/users/auth-with-password \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "secret"}'
```

**Response:**
```json
{
  "token": "eyJhbGc...",
  "record": {"id": "usr_1", "email": "admin@example.com"}
}
```

## Refresh

**POST** `/api/collections/users/auth-refresh`

```bash
curl -X POST http://localhost:8090/api/collections/users/auth-refresh \
  -H "Authorization: Bearer TOKEN"
```

## Password Reset Request

**POST** `/api/collections/users/request-password-reset`

```bash
curl -X POST http://localhost:8090/api/collections/users/request-password-reset \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com"}'
```

## Password Reset Confirm

**POST** `/api/collections/users/confirm-password-reset`

```bash
curl -X POST http://localhost:8090/api/collections/users/confirm-password-reset \
  -H "Content-Type: application/json" \
  -d '{"token": "TOKEN", "password": "newpass"}'
```
