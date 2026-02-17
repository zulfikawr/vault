# REST API

Vault provides a RESTful API for all operations.

## Base URL

```
http://localhost:8090/api
```

## Authentication

Include token in Authorization header:

```bash
curl http://localhost:8090/api/collections/posts/records \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Response Format

```json
{
  "data": {...},
  "error": null
}
```

## HTTP Methods

| Method | Operation |
|--------|-----------|
| GET | List/View |
| POST | Create |
| PATCH | Update |
| DELETE | Delete |

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Server Error |

See Also: [CRUD](./crud.md), [Auth](./auth.md)
