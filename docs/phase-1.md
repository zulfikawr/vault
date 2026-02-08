# Phase 1: Core Foundation - Detailed Implementation Plan (Standard Library Focus)

## Objective
Establish the project's structural integrity using Go's standard library where possible, focusing on logging, configuration, SQLite connectivity, and a robust HTTP server with graceful shutdown.

---

## Step 1: Project Initialization & Structure
**Goal:** Set up the Go environment and directory hierarchy.

- [x] Initialize Go module: `go mod init github.com/zulfikawr/vault`
- [x] Create the directory structure:
    - `cmd/server/`
    - `internal/core/`, `internal/db/`, `internal/api/`, `internal/models/`, `internal/auth/`
    - `pkg/`
- [x] Create a skeleton `cmd/server/main.go`.

**Deliverables:**
- `go.mod` file.
- `main.go` that prints "Starting Vault..."

---

## Step 2: Logging Infrastructure
**Goal:** Implement structured logging using the standard library.

- [x] Implement `internal/core/logger.go` using `log/slog`.
- [x] Support environment-based configuration:
    - `TextHandler` for local development.
    - `JSONHandler` for production environments.
- [x] Implement a helper to set the global log level from configuration.

**Deliverables:**
- Structured logging enabled across the app.

---

## Step 3: Native Configuration System
**Goal:** Load settings from JSON and environment variables without external libs.

- [x] Define `Config` struct in `internal/core/config.go`.
- [x] Implement `LoadConfig()`:
    - Read `config.json` if it exists (`os.ReadFile` + `json.Unmarshal`).
    - Override values with environment variables (`os.Getenv`).
- [x] Provide sensible defaults for port (8090), DB path, and log level.

**Deliverables:**
- Zero-dependency configuration loader.

---

## Step 4: SQLite Integration
**Goal:** Establish a reliable connection to the embedded database.

- [x] Implement `internal/db/connection.go` using `modernc.org/sqlite` (pure Go driver).
- [x] Configure the connection via `database/sql`:
    - Execute `PRAGMA journal_mode=WAL;` and `PRAGMA busy_timeout=5000;`.
    - Set `SetMaxOpenConns(1)` (recommended for SQLite write stability) or higher for read-heavy loads.
- [x] Ensure data directory exists (`os.MkdirAll`).

**Deliverables:**
- Verified connection to `vault.db`.

---

## Step 5: HTTP Server & Native Routing
**Goal:** Start a web server using `net/http`.

- [x] Implement `internal/core/server.go`.
- [x] Use `http.NewServeMux()` for routing (leveraging Go 1.22+ method/pattern matching).
- [x] Implement `Start()` and `Shutdown()` using `http.Server` and `context` for graceful termination.
- [x] Add `GET /api/health` handler.

**Deliverables:**
- Graceful shutdown handling for SIGINT/SIGTERM.

---

## Step 6: Custom Middleware Chain
**Goal:** Implement a lightweight middleware pattern.

- [x] Define `type Middleware func(http.Handler) http.Handler` in `internal/api/middleware.go`.
- [x] Implement a `Chain(h http.Handler, m ...Middleware)` helper function.
- [x] Create native middlewares for:
    - **Logging:** Log request method, path, status, and duration using `slog`.
    - **Recovery:** Catch panics and return 500.
    - **CORS:** Basic header injection.

**Deliverables:**
- Functional middleware pipeline without external packages.

---

## Step 7: App Bootstrapping
**Goal:** Wire all components together in a clean `App` lifecycle.

- [x] Implement `App` struct in `internal/core/app.go`.
- [x] Initialize components in order: Logger -> Config -> DB -> Server.
- [x] Update `cmd/server/main.go` to call `app.Run()`.

**Deliverables:**
- Unified entry point that manages the lifecycle of all foundation components.

---

## Verification Checklist (Phase 1 Done)
- [x] `curl http://localhost:8090/api/health` returns 200 OK.
- [x] Server port changes via `VAULT_PORT` environment variable.
- [x] Database file is created and reachable.
- [x] Logs show request details (method, path, duration).
- [x] App shuts down cleanly within 5-10 seconds of receiving Ctrl+C.
