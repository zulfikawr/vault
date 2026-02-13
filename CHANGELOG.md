# Changelog

All notable changes to Vault will be documented in this file.

## [0.0.1] - 2026-02-13

### Added
- **Core Framework**
  - Embedded SQLite database with WAL mode for high-concurrency
  - Dynamic schema engine with Collections and Fields
  - Auto-migrations for schema synchronization
  - JWT-based authentication with Bcrypt password hashing
  - Rule-based authorization with record-level security
  - Real-time subscriptions via Server-Sent Events (SSE)
  - Pluggable file storage system with local filesystem driver
  - Professional admin dashboard built with Vue 3 and Tailwind CSS

- **CLI Commands**
  - `vault serve` - Start HTTP server with port and directory options
  - `vault admin create` - Create new admin user
  - `vault admin list` - List all admin users
  - `vault admin delete` - Delete admin user with confirmation
  - `vault admin reset-password` - Reset admin password
  - `vault backup create` - Create compressed backup (zip format)
  - `vault backup list` - List all backups
  - `vault backup restore` - Restore from backup with confirmation

- **Admin Dashboard**
  - Dashboard with system overview
  - Collections management (CRUD operations)
  - Records management with dynamic field types
  - System logs viewer with filtering
  - System settings configuration
  - Consistent UI layout across all screens

- **API Endpoints**
  - Authentication: Login, refresh tokens, password reset
  - Collections: List, create, view, update, delete
  - Records: CRUD operations with filtering
  - Files: Upload and serve files
  - Real-time: SSE subscriptions
  - Admin: Collections, logs, settings management
  - Health check endpoint

- **Configuration**
  - 13 configurable options (port, database path, logging, JWT, TLS, CORS, rate limiting, file upload size)
  - Config file support (config.json)
  - Environment variable overrides (VAULT_* prefix)
  - Settings API for runtime configuration

- **Logging**
  - File-based logging system
  - Configurable log levels (DEBUG, INFO, WARN, ERROR)
  - Text and JSON log formats
  - Request ID tracking across all layers

### Features
- Single lightweight binary deployment
- No external dependencies (uses Go standard library)
- Context-aware I/O operations
- Structured error handling
- Request tracing with unique IDs
- Minimal and clean codebase

### Known Limitations
- Email verification not implemented
- TLS configuration available but not fully tested
- No built-in rate limiting enforcement yet
- Single-instance deployment (no clustering)

---

[0.0.1]: https://github.com/zulfikawr/vault/releases/tag/v0.0.1
