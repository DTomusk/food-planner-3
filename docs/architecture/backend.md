# Backend Architecture

## API
The API is a Go server with a largely flat structure. It contains a cmd directory housing executables and an internal directory which has the application logic.

Internal is divided into a number of packages, and there is little nesting beyond that. A domain package, such as recipe, is a vertical slice that contains entities, repos, and services. The only layer that is intentionally centralized is delivery (GraphQL via gqlgen), so resolvers across all domains live in gql. Resolvers should remain thin and delegate business behavior to domain services.

## Authentication
Authentication and session lifecycle behavior are documented in docs/features/auth.md.

## Testing
Testing is done throughout with unit tests and integration tests. There are currently no mocks. Repo and service tests use a real DB instance in Docker.

GraphQL schema codegen command:
go run github.com/99designs/gqlgen generate

Run from api/internal/gql.

## Deployment
Run flyctl deploy in api to deploy the API to Fly.io. A deployment workflow also deploys automatically on pushes to main.