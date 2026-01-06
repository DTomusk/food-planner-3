# Copilot Instructions - Food Planner API

## Architecture Overview
This is a Go-based food planner API using clean architecture patterns:
- **Domain Layer**: [internal/recipe/recipe.go](api/internal/recipe/recipe.go) - Core business entities with validation (e.g., `ErrEmptyName`)
- **Repository Layer**: [internal/recipe/repo.go](api/internal/recipe/repo.go) - Database operations using the `DBTX` interface for testability
- **Infrastructure**: PostgreSQL with golang-migrate, containerized with Docker Compose

## Key Patterns & Conventions

### Database Interface Pattern
All repository methods accept `db.DBTX` interface instead of `*sql.DB`:
```go
func (r *Repo) CreateRecipe(recipe Recipe, ctx context.Context, db db.DBTX) error
```
This allows both `*sql.DB` and `*sql.Tx` to be used, enabling transaction support in tests.

### Testing Strategy
- **Unit tests**: Use `testutil.WithTx()` for database tests that auto-rollback
- **Integration tests**: Docker service `test_runner` runs all tests against isolated test database
- Test database runs on port 5433, production on 5432

### Configuration Management
- [internal/config/config.go](api/internal/config/config.go) loads from `.env` file locally, environment variables in Docker
- Uses `ENV=docker` to skip `.env` loading in containers
- Database URL configured via `DB_URL`/`TEST_DB_URL` environment variables

### UUID-First Design
All entities use `github.com/google/uuid` for primary keys. Example:
```go
type Recipe struct {
    ID   uuid.UUID
    Name string
}
```

## Development Workflows

### Local Development
```bash
# Start services (from api/ directory)
docker compose up -d db migrate

# Run tests
docker compose up --build test_runner
```

### Key Docker Compose Services
- `db`: Production PostgreSQL (port 5432)
- `db_test`: Test PostgreSQL (port 5433) 
- `migrate`/`migrate_test`: golang-migrate runners
- `test_runner`: Runs `go test ./internal/... -v` in isolated environment
- `api`: Main application server

### Multi-stage Dockerfile
- `builder`: Compiles Go binary
- `test`: Extends builder for test execution
- Production: Alpine-based minimal runtime

### Migration Patterns
- Use [migrations/](api/migrations/) directory with sequential numbering (`0001_`, `0002_`, etc.)
- Up/down migration pairs required
- Auto-applied via migrate containers with health check dependencies

## File Organization
- `cmd/server/`: Application entry points and Dockerfiles
- `internal/`: Private application code following domain boundaries
- `internal/testutil/`: Shared test utilities (database setup, transaction helpers)
- Repository pattern: separate `*_test.go` files for each domain package

## When Adding New Features
1. Create domain entity in `internal/{domain}/{domain}.go` with validation
2. Add repository in `internal/{domain}/repo.go` using `DBTX` interface
3. Write tests using `testutil.WithTx()` pattern
4. Create migration files for schema changes
5. Update Docker Compose if new services needed