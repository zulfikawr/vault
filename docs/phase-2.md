# Phase 2: Database Layer - Detailed Implementation Plan

## Objective
Build the dynamic schema engine that allows users to define "Collections" (tables) and "Fields" (columns) programmatically, providing a type-safe and validated CRUD API.

---

## Step 0: Foundation Testing (Carry-over from 1.5)
**Goal:** Ensure the core utilities are bulletproof.
- [ ] Write unit tests for `internal/core/config.go`.
- [ ] Write unit tests for `internal/core/errors.go` and `internal/api/response.go`.
- [ ] Write unit tests for the middleware chain and Request ID propagation.

---

## Step 1: Collection & Field Models
**Goal:** Define how Vault stores its own metadata.
- [ ] Create `internal/models/field.go`: Define field types (text, number, bool, json, date, relation) and validation rules.
- [ ] Create `internal/models/collection.go`: Define the Collection struct (Name, Type: system/auth/base, Fields, ListRule, CreateRule, etc.).
- [ ] Implement `internal/models/record.go`: Define a generic `Record` struct that wraps `map[string]any` with helper methods for type conversion (`GetString`, `GetInt`, etc.).

---

## Step 2: Schema Registry & System Collections
**Goal:** Bootstrapping the internal tables.
- [ ] Create `internal/db/schema.go`: Logic to manage system collections (e.g., `_collections` and `_fields`).
- [ ] Implement a `SchemaRegistry` that caches collection definitions in memory for fast lookups.
- [ ] Implement logic to sync the registry with the SQLite master table.

---

## Step 3: The Migration Engine (Auto-Schema)
**Goal:** Automatically turn Collection definitions into SQL tables.
- [ ] Create `internal/db/migration.go`.
- [ ] Implement a "Sync" function that:
    - Compares the `Collection` model with the actual SQLite table structure.
    - Generates and executes `CREATE TABLE`, `ALTER TABLE`, or `CREATE INDEX` statements.
- [ ] Ensure system collections are migrated first on app startup.

---

## Step 4: Generic CRUD Operations
**Goal:** Implement the logic to interact with any collection.
- [ ] Create `internal/db/executor.go` (or update `query.go`):
    - `CreateRecord(collection, data)`: Insert with auto-generated IDs and timestamps.
    - `UpdateRecord(collection, id, data)`: Partial updates.
    - `DeleteRecord(collection, id)`: Soft/Hard delete support.
    - `FindRecordById(collection, id)`: Basic retrieval.

---

## Step 5: Data Validation Engine
**Goal:** Enforce schema rules before data touches the database.
- [ ] Implement `internal/db/validator.go`.
- [ ] Add support for:
    - Required fields.
    - String length/Regex.
    - Number ranges.
    - Unique constraints.
- [ ] Ensure validation errors return the structured `VaultError` format.

---

## Step 6: Query Builder (Filters, Sorting, Pagination)
**Goal:** Handle complex data retrieval.
- [ ] Implement `internal/db/query_builder.go`.
- [ ] Implement a simple parser for filter strings (e.g., `status = 'active'`).
- [ ] Add support for `Limit`, `Offset` (Pagination), and `OrderBy`.
- [ ] Deliver a `ListRecords(collection, params)` function.

---

## Step 7: Database Hooks
**Goal:** Allow extensibility.
- [ ] Create `internal/db/hooks.go`.
- [ ] Implement `BeforeCreate`, `AfterCreate`, `BeforeUpdate`, etc.
- [ ] Wire these hooks into the CRUD executor.

---

## Verification Checklist (Phase 2 Done)
- [ ] System tables (`_collections`) are automatically created on startup.
- [ ] Creating a new Collection model via code automatically generates a SQLite table.
- [ ] Attempting to insert a Record with a missing "required" field returns a `400 Bad Request` with structured details.
- [ ] `ListRecords` correctly handles `?page=2&perPage=10`.
- [ ] All database errors are caught, logged with `request_id`, and returned as `VaultError`.
