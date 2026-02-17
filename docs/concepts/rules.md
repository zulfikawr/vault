# Authorization Rules

Rule-based access control for collections.

## Rule Syntax

Rules use simple expressions:

```
@request.auth.id = record.author_id
@request.auth.role = 'admin'
published = true
```

## Context Variables

| Variable | Description |
|----------|-------------|
| `@request.auth.id` | Authenticated user ID |
| `@request.auth.email` | User email |
| `@request.auth.role` | User role |
| `record.field` | Record field value |

## Rule Types

| Rule | When Evaluated |
|------|----------------|
| `list_rule` | Listing records |
| `view_rule` | Viewing single record |
| `create_rule` | Creating records |
| `update_rule` | Updating records |
| `delete_rule` | Deleting records |

## Examples

### Public Read, Authenticated Write
```json
{
  "list_rule": "",
  "view_rule": "",
  "create_rule": "@request.auth.id != ''",
  "update_rule": "@request.auth.id = record.author_id",
  "delete_rule": "@request.auth.role = 'admin'"
}
```

### Owner Only
```json
{
  "list_rule": "@request.auth.id = record.user_id",
  "view_rule": "@request.auth.id = record.user_id",
  "update_rule": "@request.auth.id = record.user_id",
  "delete_rule": "@request.auth.id = record.user_id"
}
```

### Admin Only
```json
{
  "list_rule": "@request.auth.role = 'admin'",
  "view_rule": "@request.auth.role = 'admin'",
  "create_rule": "@request.auth.role = 'admin'",
  "update_rule": "@request.auth.role = 'admin'",
  "delete_rule": "@request.auth.role = 'admin'"
}
```

See Also: [Collections](./collections.md)
