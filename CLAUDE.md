# Project Overview

Go REST API server using Gin framework with CSR (Controller-Service-Repository) architecture and PostgreSQL.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **Database**: PostgreSQL (via `pgx/v5`)
- **Query**: `sqlc` (SQL → Go code generation)
- **Migration**: `goose`
- **Config**: `godotenv` + `viper`

## Project Structure

```
.
├── cmd/
│   └── main.go               # Entry point
├── internal/
│   ├── handler/              # HTTP handlers (Gin context)
│   ├── service/              # Business logic
│   ├── repository/           # DB queries (sqlc generated)
│   └── model/                # Structs / domain types
├── db/
│   ├── postgres.go           # DB connection pool
│   ├── migrations/           # goose migration files
│   └── queries/              # .sql files for sqlc
├── config/
│   └── config.go             # Env var loading
├── routes/
│   └── routes.go             # Route registration
├── middleware/
│   └── auth.go               # JWT and other middleware
├── .env                      # Local env vars (gitignored)
├── sqlc.yaml
└── go.mod
```

## Architecture Rules

- **Dependency direction**: handler → service → repository. Never skip layers.
- Handler only handles HTTP (binding, response). No business logic here.
- Service contains all business logic. No direct DB calls.
- Repository only does DB queries. Returns domain models.
- Interfaces are defined in service layer, implemented in repository.

## Common Commands

```bash
# Run server
go run cmd/main.go

# Generate sqlc code from SQL queries
sqlc generate

# Create new migration
goose -dir db/migrations create <name> sql

# Run migrations
goose -dir db/migrations postgres $DATABASE_URL up

# Rollback migration
goose -dir db/migrations postgres $DATABASE_URL down

# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Build
go build -o bin/server cmd/main.go
```

## Environment Variables

```
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
JWT_SECRET=your_secret_here
GIN_MODE=debug  # or release
```

## Code Conventions

- Use `context.Context` as the first argument in all service and repository functions.
- Always handle errors explicitly. No `_` for errors.
- Return errors up the stack. Log only at the handler layer.
- Never use `gin.H{}` for responses. Always use the `Response` struct defined below.
- Repository functions return domain models, not raw DB rows.
- Use `pgx/v5` directly via sqlc-generated code. Do not use `database/sql`.

## Response Format

All API responses MUST use the unified `Response` struct. Never use `gin.H{}`.

**Struct definition** (`internal/model/response.go`):

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
c.JSON(http.StatusOK, model.Response{
    Code:    http.StatusOK,
    Success: true,
    Message: "ok",
    Data:    result,
})

// Error
c.JSON(http.StatusBadRequest, model.Response{
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

- Write all queries in `db/queries/*.sql`
- Run `sqlc generate` after any SQL change
- Never write raw SQL strings in Go code
- Query naming convention: `GetUser`, `ListUsers`, `CreateUser`, `UpdateUser`, `DeleteUser`

## Testing

- Unit test services with mocked repository interfaces
- Use `testify` for assertions
- Test files live alongside the code (`_test.go` suffix)
- IMPORTANT: Do not test repository layer with real DB unless explicitly asked

## Git Conventions

- Branch: `feature/`, `fix/`, `chore/` prefix
- Commit: imperative mood ("Add user handler", not "Added user handler")
- Do not commit `.env` file
