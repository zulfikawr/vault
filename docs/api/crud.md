# CRUD Operations

## List Records

**GET** `/api/collections/{collection}/records`

```bash
curl http://localhost:8090/api/collections/posts/records
```

**Query Parameters:**
- `page` - Page number
- `perPage` - Items per page  
- `filter` - Filter expression
- `sort` - Sort field (- for desc)

## Create Record

**POST** `/api/collections/{collection}/records`

```bash
curl -X POST http://localhost:8090/api/collections/posts/records \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello", "body": "World"}'
```

## View Record

**GET** `/api/collections/{collection}/records/{id}`

```bash
curl http://localhost:8090/api/collections/posts/records/RECORD_ID
```

## Update Record

**PATCH** `/api/collections/{collection}/records/{id}`

```bash
curl -X PATCH http://localhost:8090/api/collections/posts/records/ID \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title": "Updated"}'
```

## Delete Record

**DELETE** `/api/collections/{collection}/records/{id}`

```bash
curl -X DELETE http://localhost:8090/api/collections/posts/records/ID \
  -H "Authorization: Bearer TOKEN"
```
