# Collections

Collections are dynamic database tables in Vault that can be created and modified on-the-fly.

## What is a Collection?

A collection is like a database table but with these features:
- Created via API or CLI (no SQL needed)
- Dynamic schema (add/remove fields anytime)
- Built-in authorization rules
- Automatic timestamps (`created`, `updated`)
- File attachment support

## Collection Structure

```json
{
  "id": "col_posts",
  "name": "posts",
  "type": "base",
  "fields": [
    {"name": "title", "type": "text", "required": true},
    {"name": "body", "type": "text"},
    {"name": "published", "type": "bool"},
    {"name": "views", "type": "number"}
  ],
  "indexes": ["published"],
  "list_rule": "@request.auth.id != ''",
  "view_rule": "published = true",
  "created": "2026-02-17T12:00:00Z",
  "updated": "2026-02-17T12:00:00Z"
}
```

## Collection Types

| Type | Description | Use Case |
|------|-------------|----------|
| `base` | Regular collection | Posts, products, comments |
| `auth` | Authentication | Users, admins |
| `system` | Internal use | _collections, _audit_logs |

## Field Types

### text
String values of any length.

```json
{"name": "title", "type": "text"}
```

### number
Numeric values (integers or decimals).

```json
{"name": "price", "type": "number"}
```

### bool
Boolean values (true/false).

```json
{"name": "published", "type": "bool"}
```

### date
Date and time values (stored as ISO 8601 strings).

```json
{"name": "published_at", "type": "date"}
```

### json
JSON objects.

```json
{"name": "metadata", "type": "json"}
```

### relation
Reference to another collection (stores record ID).

```json
{"name": "author", "type": "relation", "options": {"collection": "users"}}
```

### file
File attachments.

```json
{"name": "avatar", "type": "file"}
```

## Field Constraints

### required
Field must have a value.

```bash
vault collection create --name "posts" \
  --fields "title:text:required,body:text"
```

### unique
Field values must be unique.

```bash
vault collection create --name "users" \
  --fields "email:text:required:unique"
```

## Standard Fields

Every collection automatically has:

| Field | Type | Description |
|-------|------|-------------|
| `id` | TEXT | Primary key |
| `created` | TEXT | Creation timestamp |
| `updated` | TEXT | Last update timestamp |

## Indexes

Improve query performance:

```json
{
  "name": "posts",
  "fields": [...],
  "indexes": ["published", "author_id"]
}
```

## Authorization Rules

Control access at the record level:

| Rule | Description |
|------|-------------|
| `list_rule` | Who can list records |
| `view_rule` | Who can view single records |
| `create_rule` | Who can create records |
| `update_rule` | Who can update records |
| `delete_rule` | Who can delete records |

**Example rules:**
- `""` - Public access
- `"@request.auth.id != ''"` - Any authenticated user
- `"@request.auth.id = record.author_id"` - Owner only
- `"@request.auth.role = 'admin'"` - Admin only

## Creating Collections

### Via CLI

```bash
vault collection create \
  --name "posts" \
  --fields "title:text:required,body:text,published:bool" \
  --email "admin@example.com" \
  --password "secret"
```

### Via API

```bash
curl -X POST http://localhost:8090/api/collections \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "posts",
    "fields": [
      {"name": "title", "type": "text", "required": true},
      {"name": "body", "type": "text"}
    ]
  }'
```

### Via Admin Dashboard

1. Go to Collections
2. Click "New Collection"
3. Enter name and fields
4. Click "Create"

## Best Practices

1. **Use meaningful names**: `blog_posts` not `bp`
2. **Plan your schema**: Think about queries you'll run
3. **Add indexes**: For frequently filtered fields
4. **Set authorization rules**: Don't leave collections public
5. **Use relations**: Link collections instead of duplicating data

## See Also

- [Records](./records.md) - Working with data
- [Authorization Rules](./rules.md) - Access control
- [`vault collection`](../cli/collection.md) - CLI reference
