# Backend Architecture

## API
The API is a Go server with a largely flat structure. It contains a cmd directory housing executables and an internal directory which has the application logic.

Internal is divided into a number of packages, and there is little nesting beyond that. A domain package, such as recipe, is a vertical slice that contains entities, repos, and services. The only layer that is intentionally centralized is delivery (GraphQL via gqlgen), so resolvers across all domains live in gql. Resolvers should remain thin and delegate business behavior to domain services.

## Backend Auth Responsibilities

- The backend issues and validates short-lived access JWTs.
- The backend persists refresh tokens server-side and rotates them within refresh token families.
- Authentication middleware is permissive at the HTTP layer: valid bearer tokens populate request context, while missing or invalid tokens defer enforcement to GraphQL auth checks.
- Signout and revoke behavior invalidate the refresh token family tied to the presented refresh token.
- The full end-to-end auth and session lifecycle is documented in `docs/features/auth.md`.

### Backend Auth Source-Of-Truth Files
- api/internal/gql/graph/schema/auth.graphqls
- api/internal/gql/graph/resolver/auth.resolvers.go
- api/internal/gql/graph/resolver/auth_helpers.go
- api/internal/auth/service.go
- api/internal/auth/refresh_tokens/service.go
- api/internal/auth/middleware.go

## Testing
Testing is done throughout with unit tests and integration tests. There are currently no mocks. Repo and service tests use a real DB instance in Docker.

GraphQL schema codegen command:
go run github.com/99designs/gqlgen generate

Run from api/internal/gql.

## Deployment
Run flyctl deploy in api to deploy the API to Fly.io. A deployment workflow also deploys automatically on pushes to main.