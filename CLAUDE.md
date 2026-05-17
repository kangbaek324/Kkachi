# Project Overview

KKachi is a Go monorepo containing two apps:

- **`apps/backend`** — REST API server (Gin + PostgreSQL, CSR architecture)
- **`apps/frontend`** — Terminal UI client built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

Both modules are wired together via a top-level `go.work`.

## Tech Stack

### Backend (`apps/backend`)

- **Language**: Go 1.22+
- **Framework**: Gin
- **Database**: PostgreSQL (via `pgx/v5`)
- **Query**: `sqlc` (SQL → Go code generation)
- **Migration**: `goose`
- **Config**: `godotenv`

### Frontend (`apps/frontend`)

- **Language**: Go 1.22+
- **TUI**: `github.com/charmbracelet/bubbletea`
- **Styling**: `github.com/charmbracelet/lipgloss`

## Project Structure

```
.
├── go.work                       # Workspace pinning both modules
├── CLAUDE.md
├── README.md
├── apps/
│   ├── backend/                  # Go REST API
│   │   ├── cmd/
│   │   │   └── main.go           # Entry point
│   │   ├── internal/
│   │   │   ├── common/           # Shared types (Response, errors, etc.)
│   │   │   ├── config/           # Env var loading
│   │   │   ├── middleware/       # JWT middleware
│   │   │   └── domain/
│   │   │       ├── user/         # user handler/service/routes/dto
│   │   │       ├── wallet/       # wallet handler/service/routes
│   │   │       └── currency/     # currency handler/service/routes
│   │   ├── db/
│   │   │   ├── postgres.go       # DB connection pool
│   │   │   ├── schema.sql        # Schema reference
│   │   │   ├── migrations/       # goose migrations (sqlc schema source)
│   │   │   ├── queries/          # Raw SQL for sqlc
│   │   │   └── sqlc/             # sqlc generated code (shared)
│   │   ├── routes/
│   │   │   └── routes.go         # Route registration
│   │   ├── .env                  # Backend env vars (gitignored)
│   │   ├── .air.toml
│   │   ├── sqlc.yaml
│   │   └── go.mod                # module: github.com/kangbaek324/kkachi/apps/backend
│   └── frontend/                 # Bubble Tea TUI
│       ├── cmd/
│       │   └── main.go           # Entry point
│       ├── internal/
│       │   ├── app/              # tea.Model / Update / View
│       │   └── config/           # Env var loading
│       └── go.mod                # module: github.com/kangbaek324/kkachi/apps/frontend
```

## Architecture Rules (backend)

- **Dependency direction**: handler → service → repository. Never skip layers.
- Handler only handles HTTP (binding, response). No business logic here.
- Service contains all business logic. No direct DB calls.
- Repository only does DB queries. Returns domain models.
- Interfaces are defined in service layer, implemented in repository.

## Architecture Rules (frontend)

- Bubble Tea `Model` lives in `internal/app`. Keep `Update` pure — return `tea.Cmd` for side effects (API calls, timers).
- Network calls to the backend live in their own package (to be added when first endpoint is wired) and are invoked via `tea.Cmd`, never inline in `Update`.
- Configuration (API base URL, auth token) loads in `internal/config` and is passed into the model — do not read env vars from inside the model.

## Common Commands

All commands assume you are at the repo root unless noted.

### Backend

```bash
# Run server
go run ./apps/backend/cmd

# Generate sqlc code (run from apps/backend/)
cd apps/backend && sqlc generate

# Create new migration
goose -dir apps/backend/db/migrations create <name> sql

# Run migrations
goose -dir apps/backend/db/migrations postgres $DATABASE_URL up

# Rollback migration
goose -dir apps/backend/db/migrations postgres $DATABASE_URL down

# Run backend tests
go test ./apps/backend/...

# Build backend
go build -o bin/server ./apps/backend/cmd
```

### Frontend

```bash
# Run TUI
go run ./apps/frontend/cmd

# Run frontend tests
go test ./apps/frontend/...

# Build TUI
go build -o bin/tui ./apps/frontend/cmd
```

### Workspace

```bash
# Tidy a single module
go mod tidy -C apps/backend
go mod tidy -C apps/frontend

# Build everything
go build ./...
```

## Environment Variables

### Backend (`apps/backend/.env`)

```
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
JWT_SECRET=your_secret_here
GIN_MODE=debug  # or release
API_KEY=your_koreaexim_apikey_here
```

### Frontend (env at runtime)

```
KKACHI_API_URL=http://localhost:8080  # backend base URL (defaults to localhost:8080)
```

## Code Conventions

### Backend

- Use `context.Context` as the first argument in all service and repository functions.
- Always handle errors explicitly. No `_` for errors.
- Return errors up the stack. Log only at the handler layer.
- Never use `gin.H{}` for responses. Always use the `Response` struct defined below.
- Repository functions return domain models, not raw DB rows.
- Use `pgx/v5` directly via sqlc-generated code. Do not use `database/sql`.

### Frontend

- Keep `Update` pure. Side effects must go through `tea.Cmd`.
- Use `lipgloss` for styling — do not write raw ANSI escapes.
- Prefer narrow message types (one struct per logical event) over a single big union.

## Response Format

All API responses MUST use the unified `Response` struct. Never use `gin.H{}`.

**Struct definition** (`apps/backend/internal/common/api.response.go`):

```go
type Response struct {
    Code    int    `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}
```

**Usage in handlers**:

```go
// Success
c.JSON(http.StatusOK, common.Response{
    Code:    http.StatusOK,
    Success: true,
    Message: "ok",
    Data:    result,
})

// Error
c.JSON(http.StatusBadRequest, common.Response{
    Code:    http.StatusBadRequest,
    Success: false,
    Message: "invalid request",
})
```

**JSON output shape**:

```json
// Success
{ "code": 200, "success": true, "message": "ok", "data": { ... } }

// Error (data 필드 생략됨)
{ "code": 400, "success": false, "message": "invalid request" }
```

## SQL / sqlc Rules

- Write SQL queries in `apps/backend/db/queries/<domain>.sql`
- Schema source for sqlc is `apps/backend/db/migrations/` (not `db/schema.sql`)
- Run `sqlc generate` from inside `apps/backend/` after any SQL change
- Never write raw SQL strings in Go code
- Query naming convention: `GetUser`, `ListUsers`, `CreateUser`, `UpdateUser`, `DeleteUser`
- Generated code goes into `apps/backend/db/sqlc/` (shared package, not per-domain)

## Concurrency Rules

- In Transfer, always lock the wallet row first (`FOR UPDATE OF w`) before reading the balance. This serializes concurrent transfers on the same wallet and prevents race conditions on balance reads.
- Do not use `FOR UPDATE` on the nullable side of a `LEFT JOIN` — PostgreSQL will reject it. Use `FOR UPDATE OF <table>` to target specific tables.

## Testing

- Unit test services with mocked repository interfaces
- Use `testify` for assertions
- Test files live alongside the code (`_test.go` suffix)
- IMPORTANT: Do not test repository layer with real DB unless explicitly asked

## Git Conventions

- Branch: `feature/`, `fix/`, `chore/` prefix
- Commit: imperative mood ("Add user handler", not "Added user handler")
- Do not commit `.env` files (per-app, gitignored)
