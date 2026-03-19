# FoodSmash Architecture

This file is optimized for AI/code agents. It favors exact ownership, invariants, and change guidance over narrative explanation.

## AI Doc Map

- `docs/ai/architecture.md`: cross-system boundaries, backend-centric invariants, and full-stack runtime flow.
- `docs/ai/frontend_architecture.md`: frontend runtime, ownership boundaries, and data flow.
- `docs/ai/frontend_conventions.md`: frontend folder organization, naming, import boundaries, and placement rules.
- `docs/ai/storybook.md`: Storybook setup, decorators, and story authoring conventions.
- `docs/ai/database_schema.md`: migration-derived PostgreSQL schema reference.

## System Summary

- FoodSmash is a recipe planning application.
- The frontend lives in `web/` and is built with React, TypeScript, Vite, Tailwind, React Router, React Query, and `graphql-request`.
- Storybook is used for frontend component development, documentation, and interactive UI verification.
- The backend lives in `api/` and is a Go GraphQL API built with `gqlgen`.
- Persistent storage is PostgreSQL, with schema changes managed by SQL migrations in `api/migrations/`.
- The browser talks to the API over GraphQL. The API talks to PostgreSQL directly through repository packages.

## Repository Boundaries

- `web/`: frontend app, routes, feature hooks, GraphQL operations, generated frontend GraphQL types.
- `web/.storybook/`: Storybook runtime configuration and global preview parameters.
- `api/cmd/server/`: backend bootstrap and dependency wiring.
- `api/internal/gql/`: GraphQL schema, generated server types, resolvers, and directives.
- `api/internal/recipe/`: recipe domain models, validation, repositories, and orchestration logic.
- `api/internal/ingredient/`: ingredient reference data access and sync support.
- `api/internal/user/`: user persistence and user service logic.
- `api/internal/auth/`: JWT creation/validation and auth middleware.
- `api/internal/db/`: shared DB abstractions, including transaction helpers and the `DBTX` interface.
- `api/migrations/`: ordered database migrations with up/down pairs.
- `reference/`: YAML ingredient reference data used by sync/import flows.
- `docs/ai/`: AI-targeted operational docs like this file.

## Runtime Flow

1. The browser boots the app in `web/src/main.tsx` with `BrowserRouter`, `QueryClientProvider`, and `AuthProvider`.
2. Frontend hooks call GraphQL through `web/src/lib/graphqlClient.ts`.
3. The GraphQL client reads `VITE_API_URL` and adds a bearer token when one exists.
4. The Go server starts in `api/cmd/server/main.go`, loads config, opens PostgreSQL, builds services, and mounts the GraphQL handler at `/query`.
5. `api/internal/auth/middleware.go` parses the `Authorization` header and injects JWT claims into request context when the token is valid.
6. Missing or invalid tokens are not rejected by middleware; authorization is left to GraphQL directives and resolver logic.
7. GraphQL resolvers should stay thin and delegate to domain services.
8. Domain services orchestrate validation, transactions, and repository calls.

## Authentication Contract

- Access JWTs are short-lived and signed by backend auth services.
- Frontend stores access JWT in localStorage and sends it in Authorization Bearer headers.
- Refresh tokens are opaque values stored server-side and sent to browser as HttpOnly cookies.
- Frontend GraphQL transport must use credentials include so refresh cookie is sent on refresh requests.
- GraphQL request wrapper performs single-flight refresh and retries original operation once on UNAUTHENTICATED.
- If refresh fails, original UNAUTHENTICATED error should propagate so global sign-out handling can run.
- Refresh token families are used for session lineage and replay protection.
- Each signin/signup creates a new family, enabling multiple concurrent sessions across separate devices.
- Refresh rotates within the same family and revokes previously active tokens in that family.
- Signout/revoke invalidates the family tied to the presented refresh token.
- Signout API behavior is idempotent: success means client ends signed out, even if cookie was already missing/empty/invalid.
- Auth middleware remains permissive at HTTP layer: invalid/missing bearer token does not hard-fail request before GraphQL auth checks.

### Auth Authoritative Files

- `api/internal/gql/graph/schema/auth.graphqls`
- `api/internal/gql/graph/resolver/auth.resolvers.go`
- `api/internal/gql/graph/resolver/auth_helpers.go`
- `api/internal/auth/service.go`
- `api/internal/auth/refresh_tokens/service.go`
- `api/internal/auth/middleware.go`
- `web/src/lib/auth/token.ts`
- `web/src/lib/graphqlClient.ts`
- `web/src/lib/graphqlRequest.ts`
- `web/src/lib/auth/refresh.ts`
- `web/src/lib/auth/unauthenticated.ts`
- `web/src/app/queryClient.ts`
- `web/src/app/providers/AuthProvider.tsx`

### Auth Doc Maintenance Rule

- When auth behavior changes, update both `docs/features/auth.md` and `docs/ai/architecture.md` in the same PR.
- Update `docs/architecture/backend.md` too when backend auth ownership, package boundaries, or operational responsibilities change.

## Core Domain Model

- `RecipeContainer` is the long-lived recipe record. It belongs to a user and points at the current version.
- `RecipeVersion` contains the editable recipe content: name, ingredient usages, prep time, cook time, portions, and optional source.
- `RecipeContainer.CurrentVersionID` points to the active `RecipeVersion`.
- Updating a recipe creates a new `RecipeVersion` rather than mutating the old version in place.
- Old recipe versions remain queryable and form the recipe history.
- `RecipeSource` is version-scoped, not container-scoped.
- Ingredient usages are version-scoped, not container-scoped.
- Entities use UUIDs as primary identifiers.

## Architectural Invariants

- Repository methods should accept `db.DBTX`, not concrete `*sql.DB`, so the same code works with both plain connections and transactions.
- Write orchestration belongs in the service layer, not in GraphQL resolvers.
- GraphQL resolvers are delivery adapters, not the source of business rules.
- Recipe writes are transaction-sensitive and should persist related records as one unit.
- A recipe update must create a new version and then move the container's `current_version_id` to that new version.
- Source rows in `recipe_sources` are keyed by `recipe_version_id`; each new version needs its own source row when a source is present.
- Frontend React Query keys must be unique per data shape. Recipe detail and recipe version history should not share the same key.
- Generated code should not be hand-edited when a source file exists and a generator is available.

## Generated Code And Sources Of Truth

- Backend GraphQL schema source: `api/internal/gql/graph/schema/*.graphqls`.
- Backend GraphQL generation config: `api/internal/gql/gqlgen.yml`.
- Regenerate backend GraphQL code by running `go run github.com/99designs/gqlgen generate` from `api/internal/gql`.
- Frontend GraphQL operation sources: `web/src/**/*.graphql`.
- Frontend GraphQL generation config: `web/codegen.yml`.
- Frontend generated output: `web/src/lib/graphql.generated.ts`.
- When GraphQL schema or operation shapes change, regenerate artifacts rather than patching generated files manually.

## Configuration And Infrastructure

- Backend config is loaded by `api/internal/config/config.go`.
- Local development loads `.env`; Docker skips `.env` loading when `ENV=docker` is set.
- Required backend environment variables include:
	- `DB_URL`
	- `SERVER_PORT`
	- `CORS_ALLOWED_ORIGIN`
	- `JWT_SECRET`
	- `JWT_EXPIRATION_MINUTES`
	- `INGREDIENT_DATA_FILE_PATH`
	- `INGREDIENT_UPSERT_BATCH_SIZE`
	- `RECIPE_RETENTION_DAYS`
- CORS is configured in the server bootstrap and currently allows the configured frontend origin plus `Authorization` and `Content-Type` headers.
- Docker Compose is the primary local runtime for PostgreSQL, migrations, the API, and test services.

## Testing Model

- Backend tests prefer real database integration over mocks.
- Repository and service tests should use the transaction-based test utilities so DB changes roll back automatically.
- Migration changes should be validated with database-backed tests.
- Frontend data-layer changes should be validated by regenerating GraphQL types and checking the affected hooks/components.

## Common Change Playbooks

### Add A Persisted Recipe Field

- Add a migration in `api/migrations/`.
- Update domain structs in `api/internal/recipe/`.
- Update repository scans, inserts, and selects.
- Update service validation or orchestration if needed.
- Update GraphQL schema and resolver mapping.
- Update frontend `.graphql` documents, regenerate types, and update UI mapping.

### Change Recipe Update Behavior

- Start in `api/internal/recipe/service.go`.
- Preserve versioned history unless the change explicitly removes versioning.
- Check version-scoped tables like ingredient usages and recipe sources.
- Verify that reads still return the persisted current version after the write completes.

### Add A New Domain Slice

- Create a package under `api/internal/<domain>/`.
- Put business rules in domain/service code, not resolvers.
- Use repositories for persistence and prefer `DBTX`-compatible signatures.
- Expose the domain through GraphQL schema and thin resolvers.

### Frontend UI Or Data Changes

- Start with `docs/ai/frontend_architecture.md` for flow and ownership.
- Apply placement/import rules from `docs/ai/frontend_conventions.md`.
- For story work, follow `docs/ai/storybook.md`.

## AI Working Notes

- Start by reading the domain package and its service before changing GraphQL resolvers.
- If a change affects GraphQL shape, expect coordinated updates across schema, resolver mapping, frontend operation files, and generated types.
- Prefer editing source files that feed generators rather than editing generated output directly.
- Treat recipe versions as append-only history.
- Keep container-scoped and version-scoped data separate.
- Be careful with React Query cache keys and invalidation when adding new recipe queries.
- Frontend architecture decisions belong in `docs/ai/frontend_architecture.md`.
- Frontend placement/import conventions belong in `docs/ai/frontend_conventions.md`.
- Storybook decisions belong in `docs/ai/storybook.md`.
- There is configuration for recipe retention, but the current server bootstrap passes `nil` for retention into `recipe.NewService`; verify retention wiring before relying on it.
- Auth middleware is permissive by design: invalid tokens do not fail the HTTP request early, they simply do not populate auth context.

## Authoritative Files

- `api/cmd/server/main.go`: backend composition root.
- `api/internal/config/config.go`: environment and runtime configuration.
- `api/internal/db/db.go`: DB abstractions used by repositories/services.
- `api/internal/recipe/recipe.go`: core recipe entities and validation.
- `api/internal/recipe/service.go`: recipe orchestration and transaction boundaries.
- `api/internal/auth/middleware.go`: request auth context behavior.
- `api/internal/gql/gqlgen.yml`: backend GraphQL generation.
- `web/src/main.tsx`: frontend bootstrap.
- `web/src/lib/graphqlClient.ts`: frontend GraphQL transport.
- `web/codegen.yml`: frontend GraphQL code generation.
- `docs/ai/frontend_architecture.md`: frontend runtime flow and boundaries.
- `docs/ai/frontend_conventions.md`: frontend placement and import rules.
- `docs/ai/storybook.md`: AI-focused Storybook workflow and conventions.