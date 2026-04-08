# Database Schema

This file documents the current state of the PostgreSQL schema as of migration `0021`. It is derived from all up-migrations in `api/migrations/` and is intended as a quick reference for AI agents. Do not hand-edit column definitions here — regenerate from migrations when the schema changes.

## Schemas

- `public`: application tables
- `reference`: read-mostly reference data (ingredients)

---

## public.users

Registered user accounts.

| Column          | Type           | Constraints                  |
|-----------------|----------------|------------------------------|
| `id`            | UUID           | PRIMARY KEY                  |
| `email`         | TEXT           | NOT NULL, UNIQUE             |
| `password_hash` | TEXT           | NOT NULL                     |
| `username`      | VARCHAR(255)   | NOT NULL, UNIQUE             |

---

## public.refresh_tokens

Issued refresh tokens used for session continuation and token rotation.

| Column                | Type        | Constraints                                                    |
|-----------------------|-------------|----------------------------------------------------------------|
| `id`                  | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid()                         |
| `user_id`             | UUID        | NOT NULL, FK → `users(id)` ON DELETE CASCADE                  |
| `ip_address`          | TEXT        | NOT NULL                                                       |
| `token_hash`          | TEXT        | NOT NULL                                                       |
| `expires_at`          | TIMESTAMPTZ | NOT NULL                                                       |
| `issued_at`           | TIMESTAMPTZ | NOT NULL, DEFAULT now()                                        |
| `revoked_at`          | TIMESTAMPTZ | NULLABLE                                                       |
| `family_id`           | UUID        | NOT NULL                                                       |
| `replaced_by_token_id`| UUID        | NULLABLE, FK → `refresh_tokens(id)` ON DELETE SET NULL        |

Notes:
- Constraint `refresh_token_not_expired` enforces `expires_at > now()`.
- `family_id` groups rotated tokens in the same refresh-token lineage.
- Partial index `idx_refresh_tokens_valid` covers active rows (`revoked_at IS NULL`) by `token_hash`.
- Additional indexes exist on `user_id` and `family_id`.

---

## public.uploads

Tracks temporary upload intents and claim state for media objects (for example recipe images).

| Column               | Type        | Constraints                                                |
|----------------------|-------------|------------------------------------------------------------|
| `id`                 | UUID        | PRIMARY KEY                                                |
| `owner_user_id`      | UUID        | NOT NULL, FK → `users(id)` ON DELETE CASCADE              |
| `object_key`         | TEXT        | NOT NULL, UNIQUE                                           |
| `file_name`          | TEXT        | NOT NULL                                                   |
| `file_type`          | TEXT        | NOT NULL                                                   |
| `file_size_bytes`    | BIGINT      | NOT NULL, CHECK `>= 0`                                     |
| `purpose`            | TEXT        | NOT NULL, CHECK IN (`recipe-images`)                       |
| `created_at`         | TIMESTAMPTZ | NOT NULL, DEFAULT now()                                    |
| `expires_at`         | TIMESTAMPTZ | NOT NULL, CHECK `expires_at > created_at`                  |
| `used_at`            | TIMESTAMPTZ | NULLABLE                                                   |
| `linked_entity_type` | TEXT        | NULLABLE                                                   |
| `linked_entity_id`   | UUID        | NULLABLE                                                   |

Notes:
- Constraint `uploads_link_pair_check` enforces both link fields are null or both non-null.
- `used_at` marks successful claim by a domain entity; unclaimed rows remain `NULL`.
- Index `idx_uploads_owner_user_id` supports owner-scoped reads.
- Partial index `idx_uploads_active_unlinked` covers active unclaimed rows (`used_at IS NULL`).

---

## public.recipe_containers

The long-lived record for a recipe. Belongs to a user and tracks which version is currently active.

| Column               | Type        | Constraints                                                       |
|----------------------|-------------|-------------------------------------------------------------------|
| `id`                 | UUID        | PRIMARY KEY                                                       |
| `user_id`            | UUID        | NOT NULL, FK → `users(id)` ON DELETE CASCADE                     |
| `current_version_id` | UUID        | NULLABLE, FK → `recipe_versions(id)` DEFERRABLE INITIALLY DEFERRED |
| `created_at`         | TIMESTAMPTZ | NOT NULL, DEFAULT now()                                           |
| `deleted_on`         | TIMESTAMPTZ | NULLABLE                                                          |

Notes:
- `current_version_id` is nullable during initial insert (container is created before its first version, then updated in the same transaction).
- The FK is deferrable to allow container and version to be inserted in any order within one transaction.
- `deleted_on` implements soft delete at the container level.
- Partial index `idx_recipes_created_at_id` supports newest-first pagination for non-deleted rows.

---

## public.recipe_versions

A single immutable snapshot of a recipe's content. Recipes are never mutated in place; editing creates a new version row.

| Column      | Type        | Constraints                                      |
|-------------|-------------|--------------------------------------------------|
| `id`        | UUID        | PRIMARY KEY                                      |
| `recipe_id` | UUID        | NOT NULL, FK → `recipe_containers(id)` ON DELETE CASCADE |
| `name`      | TEXT        | NOT NULL                                         |
| `prep_mins` | INTEGER     | NOT NULL, DEFAULT 0                              |
| `cook_mins` | INTEGER     | NOT NULL, DEFAULT 0                              |
| `portions`  | INTEGER     | NOT NULL, DEFAULT 1                              |
| `version`   | INTEGER     | NOT NULL, DEFAULT 1                              |
| `img_src`   | VARCHAR(255)| NULLABLE                                         |
| `created_at`| TIMESTAMPTZ | NOT NULL, DEFAULT now()                          |

Notes:
- Version numbers are incremented by the service layer, not enforced by a DB constraint.
- Treat this table as append-only history. Old rows are never updated.
- Search indexes exist on `name`: `idx_recipe_versions_name_fts_gin` (full text) and `idx_recipe_versions_name_trgm_gin` (fuzzy trigram).
- Migration `0019` enables `pg_trgm` extension (`CREATE EXTENSION IF NOT EXISTS pg_trgm`).

---

## public.recipe_sources

Optional source attribution for a recipe version (URL, book reference, or original). One row per recipe version at most.

| Column             | Type    | Constraints                                                   |
|--------------------|---------|---------------------------------------------------------------|
| `recipe_version_id`| UUID    | PRIMARY KEY, FK → `recipe_versions(id)` ON DELETE CASCADE    |
| `type`             | INT     | NOT NULL, CHECK IN (1, 2, 3)                                  |
| `url`              | TEXT    | NULLABLE                                                      |
| `book_title`       | TEXT    | NULLABLE                                                      |
| `book_page`        | INT     | NULLABLE                                                      |
| `instructions`     | TEXT    | NULLABLE                                                      |

Type enum (enforced via CHECK constraint):
- `1` = URL — `url` must be non-null; all other source fields null.
- `2` = BookReference — `book_title` and `book_page` must be non-null; all other source fields null.
- `3` = Original — `instructions` must be non-null; all other source fields null.

Notes:
- Scoped to `recipe_version_id`, not to `recipe_containers`. A new source row is required for each new version when a source is present.

---

## public.ingredient_usages

A recipe version's use of a specific ingredient with quantity and unit.

| Column              | Type          | Constraints                                                      |
|---------------------|---------------|------------------------------------------------------------------|
| `id`                | UUID          | PRIMARY KEY                                                      |
| `recipe_version_id` | UUID          | NOT NULL, FK → `recipe_versions(id)` ON DELETE CASCADE           |
| `ingredient_id`     | UUID          | NOT NULL, FK → `reference.ingredients(id)` ON DELETE CASCADE     |
| `quantity`          | NUMERIC(10,2) | NOT NULL                                                         |
| `unit`              | INT           | NOT NULL                                                         |

Notes:
- Scoped to `recipe_version_id`. A new set of usage rows is inserted for each new version.
- `unit` is a Go `unit.Unit` integer enum validated in the service layer.

---

## reference.ingredients

Read-mostly reference table of known ingredients. Populated via the ingredient sync command.

| Column           | Type    | Constraints          |
|------------------|---------|----------------------|
| `id`             | UUID    | PRIMARY KEY          |
| `name`           | TEXT    | NOT NULL             |
| `preferred_unit` | INTEGER | NOT NULL             |
| `file_key`       | TEXT    | NOT NULL, UNIQUE     |
| `counter`        | TEXT    | NULLABLE             |
| `plural`         | TEXT    | NULLABLE             |
| `counter_plural` | TEXT    | NULLABLE             |

Notes:
- `file_key` is a stable identifier used to match rows during upsert syncs from `reference/ingredients.yaml`.
- `preferred_unit` is the unit that the service enforces when creating ingredient usages.

---

## Key Relationships

```
users
  ├── refresh_tokens (user_id)
  │     └── refresh_tokens (replaced_by_token_id)
  ├── uploads (owner_user_id)
  └── recipe_containers (user_id)
    └── recipe_versions (recipe_id)
      ├── recipe_sources (recipe_version_id)    [0..1 per version]
      └── ingredient_usages (recipe_version_id) [0..N per version]

uploads
  └── claimed by domain rows via (linked_entity_type, linked_entity_id)

reference.ingredients
  └── ingredient_usages (ingredient_id)
```

## Migration History Summary

| Migration | Description                                                    |
|-----------|----------------------------------------------------------------|
| 0001      | Create `recipes` table (id, name)                              |
| 0002      | Create `users`; add `user_id` to recipes                       |
| 0003      | Create `reference` schema and `reference.ingredients`          |
| 0004      | Create `ingredient_usages`                                     |
| 0005      | Add `prep_mins`, `cook_mins`, `portions` to recipes            |
| 0006      | Create `recipe_sources`                                        |
| 0007      | Add `deleted_on` to recipes                                    |
| 0008      | Add `counter` to `reference.ingredients`                       |
| 0009      | Add `plural`, `counter_plural` to `reference.ingredients`      |
| 0010      | Add `username` to `users`                                      |
| 0011      | Rename `recipes` → `recipe_versions`                           |
| 0012      | Rename FK columns to `recipe_version_id` in usages and sources |
| 0013      | Drop `user_id` and `deleted_on` from `recipe_versions`         |
| 0014      | Create `recipe_containers`; add `recipe_id` to versions        |
| 0015      | Add `created_at` to `recipe_versions`                          |
| 0016      | Add `version` integer to `recipe_versions`                     |
| 0017      | Create `refresh_tokens`                                        |
| 0018      | Add partial index for newest non-deleted recipe pagination     |
| 0019      | Add recipe search indexes (`fts` and `trgm`) on version name   |
| 0020      | Create `uploads` table for upload lifecycle tracking           |
| 0021      | Add `img_src` to `recipe_versions`                             |
