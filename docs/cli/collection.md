# vault collection

Manage collections (database tables).

## Usage

```bash
vault collection <subcommand> [options]
```

## Subcommands

### create

Create a new collection.

```bash
vault collection create --name NAME --fields FIELDS --email EMAIL --password PASSWORD
```

**Options:**
- `--name` (required): Collection name
- `--fields` (required): Fields in `name:type` format
- `--email` (required): Admin email
- `--password` (required): Admin password

**Field Types:**
- `text` - String values
- `number` - Numeric values
- `bool` - Boolean (true/false)
- `date` - Date/time values
- `json` - JSON objects
- `relation` - Reference to another collection
- `file` - File attachments

**Examples:**
```bash
# Simple collection
vault collection create --name "posts" \
  --fields "title:text,body:text,published:bool" \
  --email "admin@example.com" --password "secret"

# With constraints
vault collection create --name "users" \
  --fields "email:text:required:unique,name:text:required" \
  --email "admin@example.com" --password "secret"
```

### list

List all collections.

```bash
vault collection list --email EMAIL --password PASSWORD
```

**Output:**
```
Collections:
┌────────────────────────┬──────┬────────┐
│ Name                   │ Type │ Fields │
├────────────────────────┼──────┼────────┤
│ posts                  │ base │      9 │
│ users                  │ auth │      4 │
│ certificates           │ base │      9 │
└────────────────────────┴──────┴────────┘
```

### get

Get collection details.

```bash
vault collection get --name NAME --email EMAIL --password PASSWORD
```

**Output:**
```
Collection: posts
Type: base
Created: 2026-02-17T12:00:00Z
Updated: 2026-02-17T12:00:00Z

Fields:
┌──────────┬──────┬──────────┬────────┐
│ Name     │ Type │ Required │ Unique │
├──────────┼──────┼──────────┼────────┤
│ title    │ text │ false    │ false  │
│ body     │ text │ false    │ false  │
│ published│ bool │ false    │ false  │
└──────────┴──────┴──────────┴────────┘
```

### delete

Delete a collection.

```bash
vault collection delete --name NAME --email EMAIL --password PASSWORD [--force]
```

**Options:**
- `--name` (required): Collection name
- `--email` (required): Admin email
- `--password` (required): Admin password
- `--force`: Skip confirmation

**Example:**
```bash
vault collection delete --name "old-posts" \
  --email "admin@example.com" --password "secret" --force
```

## Collection Types

| Type | Description | Example |
|------|-------------|---------|
| `base` | Regular collection | posts, pages |
| `auth` | Authentication collection | users |
| `system` | System collection | _collections, _audit_logs |

## Authorization Rules

Collections support rule-based access control:

- `list_rule` - Who can list records
- `view_rule` - Who can view single records
- `create_rule` - Who can create records
- `update_rule` - Who can update records
- `delete_rule` - Who can delete records

**Example rules:**
- `""` (empty) - Public access
- `"@request.auth.id != ''"` - Any authenticated user
- `"@request.auth.id = record.author_id"` - Owner only
- `"@request.auth.role = 'admin'"` - Admin only

## Troubleshooting

### "Collection already exists"

Collection names must be unique. Delete first or use different name.

### "Invalid field type"

Supported types: text, number, bool, date, json, relation, file

### "Authentication failed"

Check email and password:
```bash
vault admin list
```

## See Also

- [`vault migrate`](./migrate.md) - Database migrations
- [Authorization Rules](../concepts/rules.md) - Access control
- [API CRUD](../api/crud.md) - Record operations
