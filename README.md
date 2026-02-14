# Vault

Vault is a self-contained, batteries-included backend framework written in Go. It provides a dynamic database schema engine, robust authentication, real-time subscriptions, and a professional administrative dashboard—all delivered as a single, lightweight binary.

## ✨ Features

- **Embedded SQLite**: Pure-Go SQLite implementation (`modernc.org/sqlite`) with WAL mode enabled for high-concurrency performance.
- **Dynamic Schema Engine**: Create and modify "Collections" (tables) and "Fields" (columns) on the fly via the API or Admin UI.
- **Auto-Migrations**: The framework automatically handles SQLite table creation and schema synchronization.
- **Identity & Auth**: Full JWT-based authentication with Bcrypt password hashing, session refresh tokens, and a protected middleware chain.
- **Rule-Based Authorization**: Fine-grained, record-level security using simple string expressions (e.g., `id = @request.auth.id`).
- **Real-time Subscriptions**: Instant event broadcasting using Server-Sent Events (SSE).
- **File Storage**: Pluggable storage system with a built-in local filesystem driver and multipart upload support.
- **Embedded Admin Dashboard**: A professional, Gruvbox-themed management interface built with Vue 3 and Vite, embedded directly into the binary.
- **Developer CLI**: Simple commands to serve the API, manage admins, and handle migrations.

## 🚀 Quick Start

### 1. Install Vault
**Linux/macOS (one-line installation):**
```bash
curl -fsSL https://raw.githubusercontent.com/zulfikawr/vault/main/install.sh | bash
```

**Or download from [GitHub Releases](https://github.com/zulfikawr/vault/releases)**

### 2. Initialize Vault
```bash
vault init --email "email@example.com" --password "yourpassword" --username "yourusername"
```

### 3. Start the server
```bash
vault serve
```

Visit `http://localhost:8090/_/` to access the Admin Dashboard.

## 🛠 CLI Usage

### Initialization
- `vault init [--email EMAIL] [--username USERNAME] [--password PASSWORD]` - Initialize new Vault project
  - `--dir DIR` - Custom data directory (default: ./vault_data)
  - `--skip-admin` - Skip admin creation
  - `--force` - Overwrite existing setup

### Server
- `vault serve [--port PORT] [--dir DIR]` - Starts the HTTP server
- `vault version` - Display current version

### Admin Management
- `vault admin create --email EMAIL --password PASSWORD --username USERNAME` - Create new admin user
- `vault admin list` - List all admin users
- `vault admin delete --email EMAIL [--force]` - Delete admin user (with confirmation)
- `vault admin reset-password --email EMAIL --password PASSWORD` - Reset admin password

### Collections
- `vault collection create --name NAME --fields FIELDS --email EMAIL --password PASSWORD` - Create new collection
- `vault collection list --email EMAIL --password PASSWORD` - List all collections
- `vault collection get --name NAME --email EMAIL --password PASSWORD` - Get collection details
- `vault collection delete --name NAME --email EMAIL --password PASSWORD [--force]` - Delete collection

### Storage Management
- `vault storage list [--path PATH] [--recursive] --email EMAIL --password PASSWORD` - List files and folders
- `vault storage create --path PATH --file FILE --email EMAIL --password PASSWORD` - Upload file to storage
- `vault storage get --path PATH --output FILE --email EMAIL --password PASSWORD [--force]` - Download file from storage
- `vault storage delete --path PATH [--recursive] [--force] --email EMAIL --password PASSWORD` - Delete file or folder

### Backup & Restore
- `vault backup create [--output FILE]` - Create backup (default: vault_backup_TIMESTAMP.zip)
- `vault backup list` - List all backups
- `vault backup restore --input FILE [--force]` - Restore from backup (with confirmation)

### Migration
- `vault migrate sync [--collection NAME] [--verbose]` - Synchronize database schema with collections
- `vault migrate status` - Show current database and collection status

## 🏗 Project Structure

```
vault/
├── cmd/vault/                  # CLI Entry point
├── internal/
│   ├── api/                    # REST API Handlers & Routing
│   │   ├── auth_handlers.go
│   │   ├── crud_handlers.go
│   │   ├── files_handlers.go
│   │   ├── storage_handlers.go
│   │   ├── admin_handlers.go
│   │   ├── middleware.go
│   │   └── router.go
│   ├── auth/                   # JWT & Password Security
│   │   ├── jwt.go
│   │   └── password.go
│   ├── cli/                    # CLI Commands
│   │   ├── admin.go
│   │   ├── backup.go
│   │   ├── collection.go
│   │   ├── init.go
│   │   ├── migrate.go
│   │   └── storage.go
│   ├── core/                   # Config, Logger, & Error System
│   │   ├── config.go
│   │   ├── logger.go
│   │   ├── file_logger.go
│   │   └── context.go
│   ├── db/                     # Schema Registry, Migration, & Executor
│   │   ├── connection.go
│   │   ├── executor.go
│   │   ├── schema.go
│   │   ├── migration.go
│   │   ├── query_builder.go
│   │   ├── validator.go
│   │   ├── audit.go
│   │   └── hooks.go
│   ├── errors/                 # Error Handling System
│   │   └── errors.go
│   ├── models/                 # Collection & Record Definitions
│   │   ├── collection.go
│   │   ├── field.go
│   │   ├── record.go
│   │   ├── user.go
│   │   └── file.go
│   ├── realtime/               # SSE Hub & Event System
│   │   ├── hub.go
│   │   └── message.go
│   ├── rules/                  # Authorization Rules Engine
│   │   └── evaluator.go
│   ├── server/                 # App Lifecycle & Server Management
│   │   ├── app.go
│   │   └── server.go
│   └── storage/                # Pluggable File Storage Drivers
│       ├── interface.go
│       └── local.go
├── ui/                         # Vue 3 Admin Dashboard
│   ├── src/
│   │   ├── components/         # Reusable Vue Components
│   │   ├── views/              # Page Components
│   │   ├── stores/             # Pinia State Management
│   │   ├── router/             # Vue Router Configuration
│   │   └── main.ts
│   ├── index.html
│   ├── vite.config.ts
│   └── tailwind.config.js
├── .github/workflows/          # CI/CD Workflows
├── go.mod & go.sum             # Go Dependencies
├── Makefile                    # Build & Development Tasks
├── CHANGELOG.md                # Version History
└── vault_data/                 # Default data directory (SQLite + Storage)
    ├── vault.db                # SQLite Database
    ├── vault.log               # Application Logs
    └── storage/                # File Storage
```