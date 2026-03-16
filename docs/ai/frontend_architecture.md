# FoodSmash Frontend Architecture

This file is optimized for AI/code agents. It documents frontend runtime flow and ownership boundaries.

## When To Read This

- You are changing frontend routing, providers, or app bootstrap.
- You are adding or modifying a feature slice under `web/src/features/`.
- You are adding frontend GraphQL operations or React Query hooks.
- You need to decide where frontend code should live.

## Frontend Summary

- Frontend root: `web/`.
- Entry point: `web/src/main.tsx`.
- App shell root: `web/src/App.tsx`.
- Routing root: `web/src/app/AppRoutes.tsx`.
- Stack: React, TypeScript, Vite, React Router, React Query, Tailwind, GraphQL Codegen, `graphql-request`.

## Frontend Ownership Map

- `web/src/app/`: application composition concerns (routes, top-level guards/layouts, providers, query client setup).
- `web/src/pages/`: route-level page components that compose feature hooks/components.
- `web/src/features/`: domain slices (auth, recipes, ingredients, users) with local GraphQL docs, hooks, mappers, and types.
- `web/src/components/`: cross-feature reusable UI/form/layout components.
- `web/src/layout/`: app-level layout wrappers used by the app shell.
- `web/src/lib/`: shared infrastructure (`graphqlClient`, auth token helpers, generated GraphQL types, shared strings).
- `web/src/test/`: shared frontend test setup/utilities.
- `web/src/stories/`: Storybook onboarding/example files. New production stories should be colocated with components.

## Runtime Flow (Frontend)

1. `web/src/main.tsx` mounts React and wraps the app with:
   - `BrowserRouter`
   - `QueryClientProvider`
   - `AuthProvider`
2. `web/src/App.tsx` wraps routes in the app shell layout.
3. `web/src/app/AppRoutes.tsx` maps URLs to page components.
4. Protected routes are gated through `ProtectedLayout`.
5. Pages orchestrate feature-level hooks/components.
6. Feature hooks call GraphQL via `web/src/lib/graphqlClient.ts` and typed documents from codegen output.

## Data Layer Flow

- GraphQL operation sources live in feature folders under `web/src/features/**/graphql/**/*.graphql`.
- Frontend codegen output is `web/src/lib/graphql.generated.ts`.
- Hooks use `useQuery` and `useMutation` from React Query with generated GraphQL types.
- `queryKey` arrays must be unique per data shape.
  - Example pattern: recipe detail `['recipe', id]` vs versions `['recipe', id, 'versions']`.
- `select` should map GraphQL response shapes into feature-facing types where needed.

## Generated Sources Of Truth (Frontend)

- GraphQL operation sources: `web/src/**/*.graphql`.
- GraphQL generation config: `web/codegen.yml`.
- Generated output: `web/src/lib/graphql.generated.ts`.
- Never hand-edit generated GraphQL output; regenerate from operation source files.

## Frontend Change Playbooks

### Add A New Feature Slice

1. Create `web/src/features/<slice>/`.
2. Add feature-local `graphql/`, `hooks/`, and `types` files as needed.
3. Keep mapping and formatting logic inside the feature (for example `mappers/`, `strings.ts`).
4. Expose stable exports via a feature `index.ts` when cross-feature usage is expected.

### Add A New Route/Page

1. Add page component in `web/src/pages/`.
2. Wire route in `web/src/app/AppRoutes.tsx`.
3. Apply `ProtectedLayout` when auth is required.
4. Keep page logic orchestration-focused; push reusable domain logic into feature hooks/components.

### Add A New Query Hook

1. Add/update operation in a feature `graphql/queries/*.graphql` file.
2. Regenerate frontend GraphQL types with `pnpm run codegen` from `web/`.
3. Implement hook in `features/<slice>/hooks/` using generated types and documents.
4. Use a stable query key and map response shape for UI consumption.

### Add A New Mutation Hook

1. Add/update operation in `graphql/mutations/*.graphql`.
2. Regenerate frontend GraphQL types.
3. Implement mutation hook with generated types.
4. Invalidate affected query keys on success.

## Authoritative Files

- `web/src/main.tsx`: frontend composition root.
- `web/src/App.tsx`: app shell root.
- `web/src/app/AppRoutes.tsx`: route map.
- `web/src/app/ProtectedLayout.tsx`: route auth gate.
- `web/src/app/queryClient.ts`: React Query client instance.
- `web/src/lib/graphqlClient.ts`: GraphQL transport and auth header behavior.
- `web/src/lib/graphql.generated.ts`: generated frontend GraphQL types/documents.
- `web/codegen.yml`: frontend GraphQL code generation config.
- `docs/ai/frontend_conventions.md`: placement and naming conventions.
- `docs/ai/storybook.md`: Storybook workflow and conventions.