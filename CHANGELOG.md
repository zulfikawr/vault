# Changelog

All notable changes to Vault will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.8.1] - 2026-02-18

### Added
- **Record Counts** - Display record counts in `vault collection list` command for better data visibility.
- **Auto-Migration** - Automatically ensure `created` and `updated` system columns exist and are populated during table migrations.
- **UI Enhancements** - Integrated custom favicon across the admin panel, sidebar, and login pages.
- **Reusable Modal** - Introduced a general-purpose `Modal` component for consistent dialogs across the application, starting with the Storage view.
- **File Management** - Added support for renaming files and folders, explicit directory creation, and batch deletion from the Storage view.
- **Batch Operations** - Implemented multi-selection and batch deletion for both collection records and storage files using a custom Checkbox component.
- **Form Polishing** - Redesigned Record New/Edit forms with field type labels, icons, and improved validation.
- **Storage Backend** - Added `Rename` and `CreateDir` methods to the Storage interface and Local implementation.
- **File Upload Options** - Added support for preserving the original filename during upload with automatic sanitization.

### Fixed
- **Record Hooks** - Fixed a bug where data modified by hooks (e.g., password hashing) was not correctly persisted to the database.
- **Type Safety** - Improved field mapping for records with safer type conversion from database values.
- **Consistency** - Ensure collections are loaded from the database during migration sync and status checks.

## [0.8.0] - 2026-02-17

### Added
- **Export Command** - Export collections and data in multiple formats
  - `vault export json` - Export all or specific collections to JSON
    - `--output FILE` - Output file path (default: `vault_export_TIMESTAMP.json`)
    - `--collection NAME` - Export specific collection only
    - `--pretty` - Pretty print JSON output
  - `vault export sql` - Export schema and data as SQL statements
    - `--output FILE` - Output SQL file path (default: `vault_export_TIMESTAMP.sql`)
  - Includes schema definitions and all records
  - Formatted table output with record counts and file sizes

- **Import Command** - Import data from external sources
  - `vault import d1` - Import from Cloudflare D1 SQL dumps
    - Automatic table creation and record insertion
    - Registers imported tables as Vault collections
    - Maps SQL types to Vault field types (TEXT→text, INTEGER→number/bool, REAL→number)
    - `--dry-run` - Validate without importing
  - `vault import json` - Import from JSON files
    - Supports Vault export format (with schema)
    - Supports simple records array format
    - `--collection NAME` - Target collection for simple imports
  - `vault import sql` - Import from generic SQL files
    - Supports CREATE TABLE and INSERT statements
  - Comprehensive error handling with structured error codes

- **Comprehensive Documentation** - Complete documentation in `docs/` directory
  - **Getting Started**: Introduction, Quick Start, Installation guides
  - **CLI Reference**: Complete documentation for all commands
    - `vault init`, `vault serve`, `vault admin`
    - `vault collection`, `vault storage`, `vault backup`
    - `vault migrate`, `vault export`, `vault import`
  - **Core Concepts**: Collections, Records, Authentication, Authorization Rules, Storage
  - **API Reference**: REST API, Authentication, CRUD, Files, Real-time
  - **Migration Guides**: From Cloudflare D1, From R2 Storage, SQL Import/Export
  - **Troubleshooting**: Common Issues, Error Codes, FAQ
  - All documentation includes working examples and best practices

## [0.7.0] - 2026-02-16

### Added
- **SQL Editor** - New interactive admin view to execute raw SQL queries.
  - **Backend**: Added `SqlService` and admin API endpoint for database query execution.
  - **UI**: Added a 2-pane SQL editor with Table and JSON result views.

## [0.6.0] - 2026-02-14

### Added
- **Storage CLI Command** - New `vault storage` command for managing files from the command line
  - `vault storage list [--path PATH] [--recursive]` - List files and folders in storage with formatted table output
  - `vault storage create --path PATH --file FILE` - Upload files to storage with automatic directory creation
  - `vault storage get --path PATH --output FILE` - Download files from storage with overwrite protection
  - `vault storage delete --path PATH [--recursive] [--force]` - Delete files or directories with confirmation prompts
  - All storage commands require admin authentication with `--email` and `--password` flags
  - Includes path validation to prevent directory traversal attacks
  - Human-readable file sizes (B, KB, MB, GB) and relative timestamps in output
  - Intelligent table truncation for long filenames (max 40 characters with "..." suffix)
  - Recursive directory listing and deletion with file/folder counting
  - Confirmation prompts for destructive operations (can be skipped with `--force`)

## [0.5.2] - 2026-02-14

### Fixed
- **Server Startup Error Handling** - Improved port availability check before starting the server
  - The server now verifies the port is available before printing the "Vault server started successfully" message.
  - Added professional error messages when the port is already in use by another process.
  - Previously, the CLI would misleadingly claim the server started successfully even if it failed to bind to the port.

- **Collection Delete Functionality** - Fixed delete button in collections table not triggering modal
  - Added missing DELETE endpoint `/api/admin/collections/{name}` for removing collections
  - Implemented proper routing for collection deletion with rate limiting
  - Added DeleteCollection method to AdminHandler with proper error handling
  - Added RemoveCollection method to SchemaRegistry to remove collections from memory
  - Added DropCollection method to MigrationEngine to properly remove tables and definitions
  - Updated UI to properly trigger confirmation modal when delete button is clicked
  - Added proper transaction handling to ensure data consistency during deletion

### Changed
- **Collection Sorting** - Collections in the sidebar now display in chronological order
  - The `/api/admin/collections` endpoint now returns collections sorted by creation date (most recent first)
  - The "Recent Collections" section in the sidebar now accurately shows the most recently created collections
  - Implemented proper timestamp parsing and comparison using RFC3339 format
  - Added fallback string comparison for edge cases where timestamp parsing fails

## [0.5.1] - 2026-02-14

### Fixed
- **CLI Help Commands** - All subcommands now properly handle `-h` and `--help` flags
  - `vault collection -h/--help` - Displays collection command usage
  - `vault admin -h/--help` - Displays admin command usage  
  - `vault backup -h/--help` - Displays backup command usage
  - `vault migrate -h/--help` - Displays migrate command usage
  - `vault init -h/--help` - Displays init command usage in consistent format
  - Previously these showed "unknown subcommand" errors instead of help text

### Changed
- **Consistent Help Format** - Standardized help output format across all commands
  - All commands now follow consistent "Usage: vault <command> [options|<subcommand> [options]]" pattern
  - All subcommands now display in consistent "Subcommands:" format
  - Init command now shows custom help format instead of Go's default flag usage
  - Main help output simplified and now includes instruction to use '-h' with commands

### Added
- **Bordered Tables** - Collection commands now display data in formatted tables with borders
  - `vault collection list` now shows collections in a bordered table format with Name, Type, and Fields columns
  - `vault collection get` now shows collection fields in a bordered table format with Name, Type, Required, and Unique columns
  - Uses Unicode box-drawing characters for clean, readable tables

## [0.5.0] - 2026-02-13

### Added
- **Collection CLI Command** - Full CRUD operations for collections via CLI
  - `vault collection create` - Create new collections with field definitions and constraints
  - `vault collection list` - List all collections with field counts
  - `vault collection get` - View full collection schema and metadata
  - `vault collection delete` - Delete collections with confirmation (--force to skip)
  - Admin authentication required for all collection operations
  - Structured logging for all collection operations
  - Transaction-based operations for data consistency

- **Data Integrity Features**
  - Transaction support for collection creation and deletion (atomic operations)
  - Field constraints: `required` and `unique` modifiers in CLI
  - Constraint enforcement in database schema (NOT NULL, UNIQUE)
  - Rollback on operation failure

- **Audit Logging System**
  - New `_audit_logs` system collection for tracking all collection changes
  - Logs collection creation and deletion with admin ID and timestamp
  - Structured JSON details for each audit event
  - Automatic audit log creation on collection operations

- **Observability & Monitoring**
  - New `/api/health/collections` endpoint for collection metrics
  - Returns collection count and health status
  - Structured logging with request IDs for all operations

- **Security Enhancements**
  - Rate limiting middleware for collection operations (10 requests/minute default)
  - CSRF token validation middleware for state-changing operations
  - Enhanced CORS middleware with origin validation
  - Admin authentication required for all collection management

## [0.4.1] - 2026-02-13

### Fixed
- Fixed storage stats API response parsing (now correctly reads nested data structure)
- Fixed storage file list API response parsing
- Removed duplicate Storage menu item from sidebar
- Fixed linting errors in storage handler (errcheck)

### Changed
- Refactored Storage.vue to use shared Table component instead of custom table markup
- Improved code reusability and consistency across admin UI
- Cleaned up TypeScript type assertions in Storage view

## [0.4.0] - 2026-02-13

### Added
- **Storage Browser UI**
  - Web interface for managing uploaded files
  - Browse storage hierarchy (collection/recordID/files)
  - Storage statistics dashboard (total files, size, collections)
  - Folder navigation with breadcrumb
  - File list with name, size, type, and modified date
  - Upload files via modal with collection and recordID selection
  - Download files
  - Delete files with confirmation dialog
  - Responsive design with Gruvbox theme
  - MIME type detection for common file types

- **Storage API Endpoints**
  - `GET /api/admin/storage` - List files and folders
  - `GET /api/admin/storage/stats` - Storage statistics
  - `DELETE /api/admin/storage` - Delete file

### Improved
- Enhanced file management capabilities
- Better visibility into storage usage
- Streamlined file upload workflow

## [0.3.0] - 2026-02-13

### Added
- **Initialization Command**
  - `vault init` - Bootstrap new Vault project with single command
    - Creates directory structure (vault_data, storage)
    - Generates secure JWT secret and config.json
    - Creates .env.example with all configuration options
    - Initializes SQLite database with system collections
    - Creates first admin user (interactive or via flags)
    - Optional `--email`, `--username`, `--password` flags for non-interactive setup
    - Optional `--dir DIR` flag for custom data directory
    - Optional `--skip-admin` flag to skip admin creation
    - Optional `--force` flag to overwrite existing setup
    - Email format and password strength validation
    - Checks for existing users before creation

### Improved
- Enhanced onboarding experience for new users
- Zero-config quick start capability
- Better project setup automation

## [0.2.0] - 2026-02-13

### Added
- **Migration Commands**
  - `vault migrate sync` - Synchronize database schema with collections
    - Optional `--collection NAME` flag to sync specific collection
    - Optional `--verbose` flag for detailed output
    - Shows migration summary with success/failure counts
  - `vault migrate status` - Display current database and collection status
    - Shows database path and total collections
    - Lists all collections with type and field count
    - Formatted table output
- **Version Command**
  - `vault version` - Display current version

### Improved
- Enhanced CLI with migration management capabilities
- Better schema synchronization control
- Detailed migration process logging

## [0.1.0] - 2026-02-13

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

[0.7.0]: https://github.com/zulfikawr/vault/releases/tag/v0.7.0
[0.6.0]: https://github.com/zulfikawr/vault/releases/tag/v0.6.0
[0.5.2]: https://github.com/zulfikawr/vault/releases/tag/v0.5.2
[0.5.1]: https://github.com/zulfikawr/vault/releases/tag/v0.5.1
[0.5.0]: https://github.com/zulfikawr/vault/releases/tag/v0.5.0
[0.4.1]: https://github.com/zulfikawr/vault/releases/tag/v0.4.1
[0.4.0]: https://github.com/zulfikawr/vault/releases/tag/v0.4.0
[0.3.0]: https://github.com/zulfikawr/vault/releases/tag/v0.3.0
[0.2.0]: https://github.com/zulfikawr/vault/releases/tag/v0.2.0
[0.1.0]: https://github.com/zulfikawr/vault/releases/tag/v0.1.0
