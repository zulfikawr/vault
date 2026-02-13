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
- **Embedded Admin Dashboard**: A professional, Monokai-themed management interface built with Vue 3 and Vite, embedded directly into the binary.
- **Developer CLI**: Simple commands to serve the API, manage admins, and handle migrations.

## 🚀 Quick Start

### 1. Build the project
Ensure you have Go and Bun (for the UI) installed.

```bash
# Build the UI
cd ui && bun install && bun x vite build && cd ..

# Build the Vault binary
go build -o vault ./cmd/vault/main.go
```

### 2. Create your first Admin
```bash
./vault admin create --email "admin@vault.local" --password "yourpassword" --username "admin"
```

### 3. Start the server
```bash
./vault serve --port 8090
```

Visit `http://localhost:8090/_/` to access the Admin Dashboard.

## 🛠 CLI Usage

### Server
- `vault serve [--port PORT] [--dir DIR]` - Starts the HTTP server

### Admin Management
- `vault admin create --email EMAIL --password PASSWORD --username USERNAME` - Create new admin user
- `vault admin list` - List all admin users
- `vault admin delete --email EMAIL [--force]` - Delete admin user (with confirmation)
- `vault admin reset-password --email EMAIL --password PASSWORD` - Reset admin password

### Backup & Restore
- `vault backup create [--output FILE]` - Create backup (default: vault_backup_TIMESTAMP.zip)
- `vault backup list` - List all backups
- `vault backup restore --input FILE [--force]` - Restore from backup (with confirmation)

## 🏗 Project Structure

```
vault/
├── cmd/vault/          # CLI Entry point
├── internal/
│   ├── api/            # REST API Handlers & Routing
│   ├── auth/           # JWT & Password Security
│   ├── core/           # Config, Logger, & Error System
│   ├── db/             # Schema Registry, Migration, & Executor
│   ├── models/         # Collection & Record Definitions
│   ├── realtime/       # SSE Hub & Event System
│   ├── server/         # App Lifecycle & Server Management
│   └── storage/        # Pluggable File Storage Drivers
├── ui/                 # Vue 3 Admin Dashboard
└── vault_data/         # Default data directory (SQLite + Storage)
```

## 📜 The Vault Standard

This project follows a strict development standard:
1. **Context-Aware**: All I/O operations (DB, Network) accept `context.Context` for tracing and cancellation.
2. **Structured Errors**: Unified `VaultError` system for consistent API feedback.
3. **Traceable**: Every request is assigned a unique `X-Request-ID` logged across all layers.
4. **Minimal Dependencies**: Priority is given to the Go Standard Library to keep the binary small and secure.
