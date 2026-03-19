# Authentication And Session Flow

This document describes how authentication works across the browser, GraphQL transport, and backend session services.

## Why This Lives In Features

Authentication is a vertical slice. It spans frontend token handling, GraphQL request behavior, backend token issuance, refresh token rotation, and signout semantics.

Architecture docs should explain ownership and boundaries. This feature doc explains what the auth system does end to end.

## Token Model

- Access token: short-lived JWT, signed by the backend, returned from signup, signin, and refresh GraphQL mutations.
- Refresh token: opaque random token, persisted server-side, sent to the browser as an HttpOnly cookie.
- Frontend JavaScript stores only the access token and cannot read refresh token cookie contents.

## Frontend Responsibilities

- The frontend stores the access token in localStorage.
- Authenticated GraphQL requests send the access token as an Authorization Bearer header.
- The GraphQL client includes credentials so the browser sends the refresh cookie when refresh is attempted.
- The request wrapper retries once after refresh on UNAUTHENTICATED and uses a single-flight guard to avoid refresh storms.
- If refresh cannot recover the session, the frontend treats the user as signed out.

## Backend Responsibilities

- The backend signs and validates access JWTs.
- The backend persists refresh tokens and uses token families to model session lineage.
- Auth middleware reads the Authorization header and injects claims into request context when the token is valid.
- Missing or invalid bearer tokens do not hard-fail the HTTP request early; GraphQL auth checks remain the enforcement point.
- Refresh rotation issues a new refresh token in the same family and revokes prior active tokens in that family.
- Signout and revoke invalidate the family associated with the presented refresh token.

## Session Model

- Each successful signup or signin creates a new refresh token family.
- Different devices or separate logins can therefore hold distinct concurrent sessions.
- Refresh requests rotate within the same family rather than creating a new family.
- Family-level revocation lets signout invalidate the current session lineage without affecting unrelated sessions.

## End-To-End Lifecycle

1. Signup or signin returns an access JWT and sets a refresh token cookie.
2. Authenticated API calls use the Authorization Bearer JWT.
3. If the access JWT expires, a refresh mutation uses the refresh cookie to obtain a new JWT.
4. The frontend stores the new JWT and retries the original request once.
5. If refresh fails, the client falls back to signed-out behavior.
6. Signout clears the cookie client-side and revokes the refresh token family server-side for that session.

## Signout Semantics

- Signout is idempotent at the API boundary.
- The resolver returns success whenever the client ends in a signed-out state, including cases where the refresh cookie is already missing, empty, or invalid after cookie clearing.
- Operational failures are the cases that should still surface as real errors.

## Source-Of-Truth Files

Backend
- api/internal/gql/graph/schema/auth.graphqls
- api/internal/gql/graph/resolver/auth.resolvers.go
- api/internal/gql/graph/resolver/auth_helpers.go
- api/internal/auth/service.go
- api/internal/auth/refresh_tokens/service.go
- api/internal/auth/middleware.go

Frontend
- web/src/lib/auth/token.ts
- web/src/lib/auth/refresh.ts
- web/src/lib/auth/unauthenticated.ts
- web/src/lib/graphqlClient.ts
- web/src/lib/graphqlRequest.ts
- web/src/app/queryClient.ts
- web/src/app/providers/AuthProvider.tsx

## Related Docs

- `docs/architecture/backend.md` for backend ownership and package boundaries.
- `docs/ai/architecture.md` for AI-facing auth invariants and maintenance rules.