# Phase 3: Authentication & Identity - COMPLETED (Steps 1-5)

## Objective
Implement a robust, secure authentication system. This includes a specialized "Auth" collection type, JWT-based session management, password hashing, and protected route middleware.

---

## Step 1: Auth Collection & User Model [DONE]
- [x] Extended `internal/models/collection.go` to support `CollectionTypeAuth`.
- [x] Updated `internal/db/schema.go` to bootstrap the default `users` collection.
- [x] Implemented `internal/models/user.go` wrapper.

## Step 2: Password Security (Bcrypt) [DONE]
- [x] Created `internal/auth/password.go` with `HashPassword` and `ComparePasswords`.
- [x] Adhered to standard error responses.

## Step 3: JWT Service [DONE]
- [x] Integrated `github.com/golang-jwt/jwt/v5`.
- [x] Created `internal/auth/jwt.go` for token generation and validation.
- [x] Updated `Config` with `JWT_SECRET` and `JWT_EXPIRY`.

## Step 4: Authentication API Endpoints [DONE]
- [x] Created `internal/api/auth_handlers.go` with `Login` handler.
- [x] Implemented `BeforeCreate` hook in `App` for auto-hashing user passwords.
- [x] Added `HideField` to `Record` to prevent leaking hashed passwords in JSON.

## Step 5: Auth Middleware (Route Protection) [DONE]
- [x] Implemented `AuthMiddleware` in `internal/api/middleware.go`.
- [x] Added `WithAuth` and `GetAuth` context helpers in `internal/core/context.go`.

---

## Verification Checklist (Phase 3 Done)
- [x] User registration automatically hashes passwords via hooks.
- [x] Login returns a signed JWT containing the `request_id`.
- [x] AuthMiddleware correctly validates tokens and injects claims into context.
- [x] Passwords are hidden in response JSON via `HideField`.