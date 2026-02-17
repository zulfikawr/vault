# Introduction to Vault

Vault is a self-contained, batteries-included backend framework written in Go. It provides a dynamic database schema engine, robust authentication, real-time subscriptions, and a professional administrative dashboard—all delivered as a single, lightweight binary.

## What is Vault?

Vault is designed to be a rapid backend solution for developers who need:
- A dynamic database schema that can be modified on-the-fly
- Built-in authentication and authorization
- Real-time data subscriptions
- File storage management
- A beautiful admin dashboard
- All in a single binary with no external dependencies

## Key Features

### 🗄️ Embedded SQLite
- Pure-Go SQLite implementation (`modernc.org/sqlite`)
- WAL mode enabled for high-concurrency performance
- No external database server required

### 🔄 Dynamic Schema Engine
- Create and modify "Collections" (tables) on the fly
- Multiple field types: text, number, boolean, date, json
- API and Admin UI support for schema changes
- Automatic migrations with `vault migrate` commands

### 🔐 Identity & Auth
- JWT-based authentication
- Bcrypt password hashing
- Session refresh tokens
- Admin user management via CLI

### 🛡️ Rule-Based Authorization
- Fine-grained, record-level security
- Simple string expressions (e.g., `id = @request.auth.id`)
- Configurable per collection

### 📡 Real-time Subscriptions
- Server-Sent Events (SSE) for instant updates
- Event broadcasting to connected clients
- Perfect for live dashboards and notifications

### 📁 File Storage
- Pluggable storage system
- Built-in local filesystem driver
- Multipart upload support
- CLI management commands

### 💾 Backup & Restore
- Full database and storage backup
- Create, list, and restore backups
- ZIP-based backup format

### 🎨 Embedded Admin Dashboard
- Professional, Gruvbox-themed interface
- Built with Vue 3 and Vite
- Embedded directly into the binary
- Real-time UI updates

### 🚀 Single Binary Distribution
- Pre-built binaries for Linux, macOS, and Windows
- One-line installation script
- No runtime dependencies

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Vault Binary                         │
├─────────────────────────────────────────────────────────┤
│  CLI  │  HTTP Server  │  Admin Dashboard (Vue 3)       │
├─────────────────────────────────────────────────────────┤
│  API Layer (REST + SSE)                                 │
├─────────────────────────────────────────────────────────┤
│  Auth  │  Collections  │  Storage  │  Realtime         │
├─────────────────────────────────────────────────────────┤
│  Schema Registry  │  Migration Engine  │  Repository   │
├─────────────────────────────────────────────────────────┤
│              Embedded SQLite (vault.db)                 │
└─────────────────────────────────────────────────────────┘
```

## Use Cases

Vault is ideal for:
- **Prototyping**: Quickly spin up a backend with dynamic schema
- **Internal Tools**: Admin dashboards and management interfaces
- **Small to Medium Apps**: Full-featured backend without complexity
- **Edge Deployments**: Single binary, minimal resources
- **Offline-First Apps**: Embedded database, no external dependencies

## When NOT to Use Vault

Consider other solutions if you need:
- Horizontal scaling across multiple nodes
- Complex relational queries with joins
- High-write throughput (millions of writes/second)
- Distributed transactions
- Advanced database features (stored procedures, triggers)

## Comparison with Alternatives

| Feature | Vault | PocketBase | Directus | Firebase |
|---------|-------|------------|----------|----------|
| Language | Go | Go | Node.js | Proprietary |
| Database | SQLite | SQLite | PostgreSQL | Firestore |
| Real-time | SSE | SSE | WebSocket | WebSocket |
| Auth | Built-in | Built-in | Built-in | Built-in |
| Storage | Local/Pluggable | Local | Cloud providers | Cloud Storage |
| Dashboard | Embedded | Embedded | Embedded | Web Console |
| Self-hosted | ✅ | ✅ | ✅ | ❌ |
| Single Binary | ✅ | ✅ | ❌ | ❌ |

## Getting Started

See the [Quick Start Guide](./quickstart.md) to get up and running in 5 minutes.

## Community & Support

- **GitHub**: https://github.com/zulfikawr/vault
- **Issues**: https://github.com/zulfikawr/vault/issues
- **Discussions**: https://github.com/zulfikawr/vault/discussions

## License

Vault is open-source software licensed under the MIT license.
