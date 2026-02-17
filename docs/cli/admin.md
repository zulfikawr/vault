# vault admin

Manage admin users.

## Usage

```bash
vault admin <subcommand> [options]
```

## Subcommands

### create

Create a new admin user.

```bash
vault admin create --email EMAIL --username USERNAME --password PASSWORD
```

**Options:**
- `--email` (required): Admin email
- `--username` (required): Admin username  
- `--password` (required): Admin password

**Example:**
```bash
vault admin create --email "admin@example.com" --username "admin" --password "securepass123"
```

### list

List all admin users.

```bash
vault admin list
```

**Output:**
```
Admin Users:
┌────┬─────────────────────┬──────────┬──────────────────────┐
│ ID │ Email               │ Username │ Created              │
├────┼─────────────────────┼──────────┼──────────────────────┤
│ 1  │ admin@example.com   │ admin    │ 2026-02-17 12:00:00  │
└────┴─────────────────────┴──────────┴──────────────────────┘
```

### delete

Delete an admin user.

```bash
vault admin delete --email EMAIL [--force]
```

**Options:**
- `--email` (required): Email to delete
- `--force`: Skip confirmation

**Example:**
```bash
vault admin delete --email "old@example.com" --force
```

### reset-password

Reset admin password.

```bash
vault admin reset-password --email EMAIL --password PASSWORD
```

**Options:**
- `--email` (required): Admin email
- `--password` (required): New password

**Example:**
```bash
vault admin reset-password --email "admin@example.com" --password "newpassword123"
```

## Security Notes

- Passwords are hashed with bcrypt
- Minimum password length: 8 characters
- Email must be unique
- Username must be unique

## Troubleshooting

### "User already exists"

Email or username is already taken. Use different credentials.

### "User not found"

Check the email address is correct:
```bash
vault admin list
```

## See Also

- [`vault init`](./init.md) - Initialize with first admin
- [Authentication](../concepts/auth.md) - Auth system
