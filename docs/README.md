# Vault Documentation

Welcome to the Vault documentation. Vault is a self-contained, batteries-included backend framework written in Go.

## Table of Contents

### Getting Started
- [Introduction](./introduction.md) - What is Vault and why use it
- [Quick Start](./quickstart.md) - Get up and running in 5 minutes
- [Installation](./installation.md) - Installation methods and requirements

### Core Concepts
- [Collections](./concepts/collections.md) - Dynamic database schema
- [Records](./concepts/records.md) - Data management
- [Authentication](./concepts/auth.md) - JWT-based auth system
- [Authorization Rules](./concepts/rules.md) - Row-level security
- [Storage](./concepts/storage.md) - File storage system

### CLI Reference
- [Overview](./cli/overview.md) - CLI structure and usage
- [vault init](./cli/init.md) - Initialize new project
- [vault serve](./cli/serve.md) - Start the server
- [vault admin](./cli/admin.md) - Manage admin users
- [vault collection](./cli/collection.md) - Manage collections
- [vault storage](./cli/storage.md) - Manage file storage
- [vault export](./cli/export.md) - Export data
- [vault import](./cli/import.md) - Import data
- [vault backup](./cli/backup.md) - Backup and restore
- [vault migrate](./cli/migrate.md) - Database migrations

### API Reference
- [REST API](./api/rest.md) - RESTful API endpoints
- [Authentication](./api/auth.md) - Auth endpoints
- [CRUD Operations](./api/crud.md) - Create, read, update, delete
- [File Upload](./api/files.md) - File handling
- [Real-time](./api/realtime.md) - SSE subscriptions

### Migration Guides
- [From Cloudflare D1](./migration/from-d1.md) - Migrate from Wrangler D1
- [From R2 Storage](./migration/from-r2.md) - Migrate from R2 buckets
- [SQL Import/Export](./migration/sql.md) - SQL data migration

### Advanced Topics
- [Configuration](./advanced/config.md) - Configuration options
- [Security](./advanced/security.md) - Security best practices
- [Performance](./advanced/performance.md) - Optimization tips
- [Deployment](./advanced/deployment.md) - Production deployment

### Troubleshooting
- [Common Issues](./troubleshooting/common.md) - Frequently encountered problems
- [Error Codes](./troubleshooting/errors.md) - Error code reference
- [FAQ](./troubleshooting/faq.md) - Frequently asked questions

---

## Quick Links

- **GitHub Repository**: https://github.com/zulfikawr/vault
- **Admin Dashboard**: `http://localhost:8090/_/` (default)
- **API Endpoint**: `http://localhost:8090/api/` (default)

## Getting Help

1. Check the [FAQ](./troubleshooting/faq.md)
2. Search [existing issues](https://github.com/zulfikawr/vault/issues)
3. Create a new issue with detailed information

## Version Information

Current Version: 0.7.0

See [CHANGELOG.md](../CHANGELOG.md) for version history.
