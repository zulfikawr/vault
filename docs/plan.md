# Vault - Go Backend Framework (PocketBase Alternative)

## Project Overview
Build a self-contained, batteries-included backend framework in Go that provides database, authentication, real-time subscriptions, and admin UI out of the box.

## Core Architecture

### 1. Foundation Layer
- **Embedded SQLite Database**
  - Use `modernc.org/sqlite` (pure Go) or `mattn/go-sqlite3` (CGO)
  - Connection pooling and WAL mode
  - Migration system with version tracking
  - Schema introspection and validation

- **HTTP Server & Router**
  - Custom router or use `chi` for lightweight routing
  - Middleware chain (CORS, logging, recovery, rate limiting)
  - Request context management
  - Graceful shutdown handling

- **Configuration System**
  - YAML/JSON config file support
  - Environment variable overrides
  - Hot-reload capability for development
  - Validation on startup

### 2. Database Layer

- **ORM/Query Builder**
  - Schema definition DSL
  - CRUD operations with type safety
  - Query builder with filtering, sorting, pagination
  - Relation handling (one-to-many, many-to-many)
  - Hooks system (before/after create, update, delete)

- **Collections System**
  - Dynamic schema creation
  - Field types: text, number, bool, email, url, date, json, file, relation
  - Field validation rules
  - Indexes management
  - System collections vs user collections

- **Migration Engine**
  - Auto-migration from schema changes
  - Manual migration support
  - Rollback capability
  - Migration history tracking

### 3. Authentication & Authorization

- **Auth System**
  - Email/password authentication
  - OAuth2 providers (Google, GitHub, etc.)
  - JWT token generation and validation
  - Refresh token mechanism
  - Password reset flow
  - Email verification

- **Authorization**
  - Role-based access control (RBAC)
  - Collection-level permissions
  - Record-level rules (filter expressions)
  - API rules: list, view, create, update, delete
  - Admin vs regular user separation

- **Session Management**
  - Token storage and invalidation
  - Multi-device session tracking
  - Session expiry and cleanup

### 4. API Layer

- **REST API**
  - Auto-generated CRUD endpoints per collection
  - Filtering: `?filter=(status='active' && age>18)`
  - Sorting: `?sort=-created,+name`
  - Pagination: `?page=1&perPage=20`
  - Field expansion: `?expand=author,comments`
  - Batch operations support

- **Real-time Subscriptions**
  - WebSocket or SSE implementation
  - Subscribe to collection changes
  - Event types: create, update, delete
  - Filter subscriptions by query
  - Connection management and cleanup

- **File Storage**
  - Multiple storage backends (local, S3, R2, MinIO)
  - Pluggable storage interface
  - Image resizing on-the-fly with caching
  - File validation (size, type, dimensions)
  - Secure file serving with signed URLs/tokens
  - Thumbnail generation
  - Multiple file uploads per record
  - Direct upload support (presigned URLs)
  - CDN integration ready

### 5. Admin Dashboard

- **Backend API**
  - Collection management (CRUD)
  - Schema editor endpoints
  - User management
  - Settings management
  - Logs viewer API
  - Backup/restore endpoints

- **Frontend (Embedded)**
  - SPA built with Svelte/React (embedded via `embed` package)
  - Collection browser and editor
  - Schema designer UI
  - User management interface
  - API rules editor
  - Real-time logs viewer
  - Settings panel

### 6. Plugin System

- **Hook System**
  - Request hooks (before/after)
  - Database hooks (before/after CRUD)
  - Auth hooks (login, register, etc.)
  - Custom route registration
  - Middleware injection

- **Extension API**
  - Plugin interface definition
  - Plugin discovery and loading
  - Plugin configuration
  - Plugin lifecycle management

### 7. Developer Experience

- **CLI Tool**
  - `serve` - Start server
  - `migrate` - Run migrations
  - `admin` - Create admin user
  - `backup` - Database backup
  - `restore` - Database restore
  - `generate` - Code generation

- **SDK Generation**
  - Auto-generate TypeScript/JavaScript SDK
  - Type definitions from schema
  - API client with auth handling

- **Development Tools**
  - Hot reload for Go code changes
  - Request logging and debugging
  - Schema validation on startup
  - Health check endpoints

## Implementation Phases

### Phase 1: Core Foundation (Weeks 1-2)
- [ ] Project structure and module setup
- [ ] HTTP server with router
- [ ] SQLite integration and connection management
- [ ] Basic configuration system
- [ ] Logging infrastructure
- [ ] Middleware framework

### Phase 2: Database Layer (Weeks 3-4)
- [ ] Schema definition system
- [ ] Collections management
- [ ] Query builder implementation
- [ ] CRUD operations
- [ ] Migration engine
- [ ] Field validation system

### Phase 3: Authentication (Week 5)
- [ ] User model and auth collection
- [ ] JWT token generation/validation
- [ ] Email/password auth endpoints
- [ ] Password hashing (bcrypt)
- [ ] Token refresh mechanism
- [ ] Password reset flow

### Phase 4: REST API (Weeks 6-7)
- [ ] Auto-generated CRUD endpoints
- [ ] Query filtering parser
- [ ] Sorting and pagination
- [ ] Relation expansion
- [ ] Request validation
- [ ] Error handling standardization

### Phase 5: Authorization (Week 8)
- [ ] Permission rules engine
- [ ] Collection-level rules
- [ ] Record-level filtering
- [ ] Admin middleware
- [ ] API rule evaluation

### Phase 6: File Storage (Week 9)
- [ ] Storage interface design
- [ ] Local filesystem implementation
- [ ] S3/R2 adapter (using AWS SDK)
- [ ] File upload handling (multipart)
- [ ] File serving with auth/signed URLs
- [ ] Image resizing (using `imaging` library)
- [ ] Thumbnail generation and caching
- [ ] File validation (MIME type, size, dimensions)
- [ ] Direct upload with presigned URLs
- [ ] File metadata storage in DB
- [ ] Cleanup orphaned files

### Phase 7: Real-time (Week 10)
- [ ] WebSocket server setup
- [ ] Subscription management
- [ ] Event broadcasting system
- [ ] Connection authentication
- [ ] Filter-based subscriptions

### Phase 8: Admin Dashboard Backend (Week 11)
- [ ] Collection management API
- [ ] Schema editor endpoints
- [ ] User management API
- [ ] Settings API
- [ ] Logs API
- [ ] Backup/restore API

### Phase 9: Admin Dashboard Frontend (Weeks 12-13)
- [ ] UI framework setup (Svelte/React)
- [ ] Collection browser
- [ ] Schema designer
- [ ] User management UI
- [ ] API rules editor
- [ ] Build and embed process

### Phase 10: Plugin System (Week 14)
- [ ] Hook system implementation
- [ ] Plugin interface design
- [ ] Plugin loader
- [ ] Example plugins
- [ ] Plugin documentation

### Phase 11: CLI & Developer Tools (Week 15)
- [ ] CLI framework
- [ ] Serve command
- [ ] Migration commands
- [ ] Admin creation command
- [ ] Backup/restore commands
- [ ] Hot reload for development

### Phase 12: SDK Generation (Week 16)
- [ ] Schema to TypeScript types
- [ ] API client generator
- [ ] Auth handling in SDK
- [ ] Real-time client
- [ ] SDK documentation

### Phase 13: Testing & Documentation (Weeks 17-18)
- [ ] Unit tests for core components
- [ ] Integration tests
- [ ] API endpoint tests
- [ ] Documentation site
- [ ] Example projects
- [ ] Migration guide from PocketBase

### Phase 14: Performance & Security (Week 19)
- [ ] Query optimization
- [ ] Connection pooling tuning
- [ ] Rate limiting
- [ ] Security audit
- [ ] SQL injection prevention
- [ ] XSS protection

### Phase 15: Polish & Release (Week 20)
- [ ] Bug fixes
- [ ] Performance benchmarks
- [ ] Docker image
- [ ] Release binaries for multiple platforms
- [ ] Changelog and versioning

## Project Structure

```
project/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── core/
│   │   ├── app.go               # Main app struct
│   │   ├── config.go            # Configuration
│   │   └── server.go            # HTTP server
│   ├── db/
│   │   ├── connection.go        # DB connection
│   │   ├── schema.go            # Schema management
│   │   ├── query.go             # Query builder
│   │   ├── migration.go         # Migrations
│   │   └── hooks.go             # DB hooks
│   ├── models/
│   │   ├── collection.go        # Collection model
│   │   ├── record.go            # Record model
│   │   ├── user.go              # User model
│   │   └── field.go             # Field types
│   ├── auth/
│   │   ├── jwt.go               # JWT handling
│   │   ├── oauth.go             # OAuth providers
│   │   ├── password.go          # Password utils
│   │   └── middleware.go        # Auth middleware
│   ├── api/
│   │   ├── router.go            # Route setup
│   │   ├── crud.go              # CRUD handlers
│   │   ├── auth.go              # Auth endpoints
│   │   ├── files.go             # File endpoints
│   │   └── admin.go             # Admin endpoints
│   ├── realtime/
│   │   ├── websocket.go         # WS server
│   │   ├── subscription.go      # Subscription manager
│   │   └── events.go            # Event system
│   ├── storage/
│   │   ├── interface.go         # Storage interface
│   │   ├── local.go             # Local filesystem
│   │   ├── s3.go                # S3/R2/MinIO adapter
│   │   ├── image.go             # Image processing
│   │   ├── cache.go             # Thumbnail cache
│   │   └── validator.go         # File validation
│   ├── rules/
│   │   ├── parser.go            # Rule parser
│   │   ├── evaluator.go         # Rule evaluation
│   │   └── context.go           # Evaluation context
│   └── plugins/
│       ├── loader.go            # Plugin loader
│       ├── hooks.go             # Hook system
│       └── interface.go         # Plugin interface
├── pkg/
│   └── client/                  # Go SDK
├── ui/
│   ├── src/                     # Admin dashboard source
│   └── dist/                    # Built assets
├── migrations/
│   └── system/                  # System migrations
├── examples/
│   └── blog/                    # Example project
├── docs/
│   └── api/                     # API documentation
└── scripts/
    └── build.sh                 # Build scripts
```

## Key Differentiators from PocketBase

1. **Plugin Architecture**: More extensible plugin system from the start
2. **Multi-Database Support**: Design for future PostgreSQL/MySQL support
3. **Advanced Query Language**: More powerful filtering and aggregation
4. **Storage Flexibility**: Multiple backends (local, S3, R2, custom) with interface
5. **Built-in Caching**: Redis integration for performance + thumbnail caching
6. **GraphQL Option**: Optional GraphQL API alongside REST
7. **Workflow Engine**: Built-in background job processing
8. **Multi-tenancy**: Native support for multi-tenant applications
9. **Audit Logging**: Comprehensive audit trail system
10. **Direct Uploads**: Presigned URL support for client-side uploads

## Success Metrics

- Single binary deployment
- < 20MB binary size
- < 10ms average API response time
- Support 10k+ concurrent connections
- < 5 minute setup time for new projects
- Comprehensive test coverage (>80%)
- Clear documentation with examples

## Risk Mitigation

- **Performance**: Benchmark early and often, optimize hot paths
- **Security**: Regular security audits, follow OWASP guidelines
- **Compatibility**: Extensive testing across platforms
- **Breaking Changes**: Semantic versioning, migration guides
- **Community**: Clear contribution guidelines, responsive to issues
