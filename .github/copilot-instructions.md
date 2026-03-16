# Copilot Instructions - FoodSmash

## Reference Docs

Before making changes, read the relevant docs in `docs/ai/`:
- `docs/ai/architecture.md` — system boundaries, domain model, runtime flow, change playbooks, and architectural invariants.
- `docs/ai/database_schema.md` — current PostgreSQL schema derived from migrations. Use this to understand table structure without reading all migration files.
- `docs/ai/storybook.md` — Storybook structure, story authoring conventions, decorators, and Tailwind/app-context parity guidance.

## Constraints

### Testing
- **Never run Go tests directly** (e.g. `go test ./...`). All backend tests require a real PostgreSQL test database and must be run through Docker Compose: `docker compose -f api/docker-compose.test.yml up --build test_runner`. Running tests outside Docker will fail due to missing `TEST_DB_URL`.
- **Never use mocks** for database tests. The project uses real DB integration tests with transaction rollback via `testutil.WithTx()`.

### Database
- **Never hand-edit generated or migration-derived files** such as `docs/ai/database_schema.md`. Update migrations instead and regenerate docs.
- **Always add up/down migration pairs** when changing the schema. Never modify existing migration files.

### GraphQL
- **Never hand-edit generated GraphQL files** (`web/src/lib/graphql.generated.ts` or the gqlgen output in `api/internal/gql/graph/`). Regenerate them from their sources using `go run github.com/99designs/gqlgen generate` (backend) or `pnpm run codegen` (frontend) after changing schema or operation files.