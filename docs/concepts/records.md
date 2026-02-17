# Records

Records are individual data entries in a collection.

## Creating Records

```bash
curl -X POST http://localhost:8090/api/collections/posts/records \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello", "body": "World"}'
```

## Reading Records

```bash
# List all
curl http://localhost:8090/api/collections/posts/records

# With filter
curl "http://localhost:8090/api/collections/posts/records?filter=published=true"

# Single record
curl http://localhost:8090/api/collections/posts/records/RECORD_ID
```

## Updating Records

```bash
curl -X PATCH http://localhost:8090/api/collections/posts/records/ID \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title": "Updated Title"}'
```

## Deleting Records

```bash
curl -X DELETE http://localhost:8090/api/collections/posts/records/ID \
  -H "Authorization: Bearer TOKEN"
```

## Query Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `page` | Page number | `?page=2` |
| `perPage` | Items per page | `?perPage=20` |
| `filter` | Filter expression | `?filter=published=true` |
| `sort` | Sort field | `?sort=-created` |

## Filters

- Equality: `field=value`
- Comparison: `field>10`, `field<100`
- Like: `field~pattern`
- And: `field1=val1,field2=val2`

See Also: [API CRUD](../api/crud.md)
