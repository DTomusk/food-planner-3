# FoodSmash Frontend Conventions

This file is optimized for AI/code agents. It defines frontend code organization and placement rules.

## When To Read This

- You are creating, moving, or renaming frontend files.
- You are unsure where a component, hook, or GraphQL file belongs.
- You are enforcing import boundaries between shared code and feature slices.

## Folder Organization Rules

- `web/src/app/`:
  - App-level composition only (routing, provider wiring, route guards, query client setup).
  - Do not place feature business logic here.
- `web/src/pages/`:
  - Route-level components.
  - Pages compose feature hooks/components; they should not become shared libraries.
- `web/src/features/<slice>/`:
  - Domain-specific code for one business area.
  - Keep feature-local GraphQL operations, hooks, types, mappers, schemas, and feature components here.
- `web/src/components/`:
  - Cross-feature reusable UI/form/layout components.
  - If only one feature uses a component, keep it inside that feature.
- `web/src/layout/`:
  - Reusable app shell/layout wrappers.
- `web/src/lib/`:
  - Shared infrastructure and integration helpers (`graphqlClient`, generated types, auth token helpers).
- `web/src/test/`:
  - Shared test utilities/setup.
- `web/src/stories/`:
  - Storybook onboarding/example files.
  - New component stories should be colocated with source components.

## File Naming And Placement

- React components: `PascalCase.tsx`.
- Hooks: `useX.ts`.
- Tests: colocate as `*.test.ts` or `*.test.tsx`.
- Stories: colocate as `Component.stories.ts` or `Component.stories.tsx`.
- Feature types: keep in `types.ts` unless split size demands a folder.
- Feature constants/labels: keep in `strings.ts` when feature-specific.
- GraphQL operations:
  - Place in `features/<slice>/graphql/queries/` or `features/<slice>/graphql/mutations/`.
  - Use descriptive operation filenames (for example `getRecipe.graphql`, `updateRecipe.graphql`).

## Import Boundary Rules

- Use the path alias `@/*` (base `web/src`) for cross-folder imports.
- Prefer relative imports for files inside the same feature slice.
- Avoid importing from another feature's internal subpaths when a public feature export exists.
- Keep generated GraphQL imports sourced from `web/src/lib/graphql.generated.ts` (directly or via `web/src/lib/index.ts`).

## React Query Conventions

- Query keys must be stable arrays and unique per data shape.
- Mutation success handlers should invalidate only affected keys.
- Prefer mapping GraphQL response types into feature-level UI types in hook `select` blocks or mappers.

## Storybook Conventions

- Use public Storybook types from `@storybook/react-vite`.
- Do not import Storybook types from `storybook/internal/*`.
- Keep story args/options aligned with real component prop options.
- Follow `docs/ai/storybook.md` for decorators, global styling parity, and story state coverage.

## PR Checklist (Frontend Structure)

1. New files are placed in the correct layer (`app`, `pages`, `features`, `components`, `lib`).
2. Feature-specific logic is not moved into shared folders prematurely.
3. Imports follow alias/relative rules and avoid cross-feature deep coupling.
4. Any new GraphQL operation was followed by `pnpm run codegen` in `web/`.
5. New reusable UI components include colocated stories.