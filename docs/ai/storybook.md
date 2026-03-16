# FoodSmash Storybook Guide (For AI Agents)

This file defines how Storybook is configured and how stories should be authored in this repository.

## Purpose

- Use Storybook to develop and verify UI components in isolation.
- Keep component docs and UI states close to source code.
- Ensure Storybook rendering matches application rendering (styling and provider context parity).

## Current Setup

- Storybook framework: `@storybook/react-vite`.
- Story discovery is configured in `web/.storybook/main.ts`:
  - `../src/**/*.mdx`
  - `../src/**/*.stories.@(js|jsx|mjs|ts|tsx)`
- Global Storybook preview config is in `web/.storybook/preview.ts`.
- Global Tailwind/app CSS for stories is imported via `web/.storybook/preview.ts` from `web/src/index.css`.
- Tailwind is configured in the Vite app plugin chain in `web/vite.config.ts` via `@tailwindcss/vite`.

## Story Authoring Conventions

- Prefer colocated stories next to components:
  - Example: `web/src/components/ui/Button.tsx` with `web/src/components/ui/Button.stories.ts`.
- Use Storybook public types, not internal APIs:
  - Use `Meta` and `StoryObj` from `@storybook/react-vite`.
  - Do not import from `storybook/internal/*`.
- Use default export metadata with `satisfies Meta<typeof Component>`.
- Prefer `tags: ['autodocs']` for component docs generation.
- Keep story titles stable and hierarchical (for example, `UI/Button`).
- Keep controls and args aligned with component prop unions/options.
  - Do not expose control options that the component cannot render.

## Required Story States

For reusable UI components, include at least:

- Primary/default state.
- Main variants (for example, semantic variants).
- Disabled state.
- Loading/busy state when supported.
- Any critical edge state (error, empty, destructive action) when applicable.

## Styling Parity Rules

- Keep `import '../src/index.css';` in `web/.storybook/preview.ts`.
- If global styles move, update Storybook preview import in the same PR.
- When adding design tokens/theme variables in app CSS, verify Storybook still renders tokens correctly.

## App Context Parity Rules

Some components depend on app providers from `web/src/main.tsx`:

- `BrowserRouter`
- `QueryClientProvider`
- `AuthProvider`

In Storybook, use decorators to provide equivalent context when needed:

- Prefer `MemoryRouter` for route context.
- Use a fresh `QueryClient` per story render to avoid cross-story cache leakage.
- Provide a controlled auth context/decorator for authenticated and anonymous story variants.

Keep decorators lightweight and composable. Avoid pulling full app bootstrapping into every story.

## Commands

Run from `web/`:

- `pnpm storybook` to run Storybook locally.
- `pnpm build-storybook` to build static Storybook output.

## Common Change Playbook

### Add A Story For A New Component

1. Create `Component.stories.ts[x]` next to `Component.tsx`.
2. Add typed `meta` using `Meta` and `StoryObj`.
3. Add default args and relevant controls.
4. Add primary, variant, and edge-state stories.
5. Run `pnpm storybook` and validate controls/actions/docs output.

### Move Existing Story To Component Folder

1. Copy story file next to component and convert imports to relative component imports.
2. Keep story title unchanged unless intentional taxonomy update is needed.
3. Delete the old duplicate story file.
4. Verify no duplicate entries appear in Storybook sidebar.

### Add Provider-Dependent Stories

1. Add or reuse preview decorators in `web/.storybook/preview.ts`.
2. Wrap story render with only the providers required by that component.
3. Avoid hidden network coupling in stories; pass explicit args/mocks through props where possible.

## AI Working Notes

- When touching Storybook files, check both `web/.storybook/main.ts` and `web/.storybook/preview.ts` first.
- If styling is missing in stories, check preview CSS import before changing component classes.
- If story types fail unexpectedly, verify imports are from `@storybook/react-vite`.
- Prefer story colocation for all new component work.