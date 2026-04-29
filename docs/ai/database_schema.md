# Database Schema

This file documents the current state of the PostgreSQL schema as of migration `0031`. It is derived from all up-migrations in `api/migrations/` and is intended as a quick reference for AI agents. Do not hand-edit column definitions here — regenerate from migrations when the schema changes.

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
| `deleted_at`         | TIMESTAMPTZ | NULLABLE                                                   |

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
| `published_at`       | TIMESTAMPTZ | NULLABLE                                                          |
| `deleted_on`         | TIMESTAMPTZ | NULLABLE                                                          |

Notes:
- `current_version_id` is nullable during initial insert (container is created before its first version, then updated in the same transaction).
- The FK is deferrable to allow container and version to be inserted in any order within one transaction.
- `deleted_on` implements soft delete at the container level.
- Partial index `idx_recipes_created_at_id` supports newest-first pagination for non-deleted rows.
- Partial index `idx_recipe_containers_public_created_at_id` supports newest-first pagination for published, non-deleted rows.

---

## public.recipe_versions

A single immutable snapshot of a recipe's content. Recipes are never mutated in place; editing creates a new version row.

| Column                 | Type         | Constraints                                      |
|------------------------|--------------|--------------------------------------------------|
| `id`                   | UUID         | PRIMARY KEY                                      |
| `recipe_id`            | UUID         | NOT NULL, FK → `recipe_containers(id)` ON DELETE CASCADE |
| `name`                 | TEXT         | NOT NULL                                         |
| `prep_mins`            | INTEGER      | NOT NULL, DEFAULT 0                              |
| `cook_mins`            | INTEGER      | NOT NULL, DEFAULT 0                              |
| `portions`             | INTEGER      | NOT NULL, DEFAULT 1                              |
| `version`              | INTEGER      | NOT NULL, DEFAULT 1                              |
| `img_src`              | VARCHAR(255) | NULLABLE                                         |
| `created_at`           | TIMESTAMPTZ  | NOT NULL, DEFAULT now()                          |
| `published_at`         | TIMESTAMPTZ  | NULLABLE                                         |
| `description`          | TEXT         | NOT NULL, DEFAULT ''                             |
| `animal_product_level` | INTEGER      | NOT NULL, DEFAULT 0                              |
| `contains_gluten`      | BOOLEAN      | NOT NULL, DEFAULT FALSE                          |

Notes:
- Version numbers are incremented by the service layer, not enforced by a DB constraint.
- Treat this table as append-only history. Old rows are never updated.
- Search indexes exist on `name`: `idx_recipe_versions_name_fts_gin` (full text) and `idx_recipe_versions_name_trgm_gin` (fuzzy trigram).
- Partial unique index `ux_recipe_versions_one_draft_per_recipe` enforces at most one draft row (`published_at IS NULL`) per `recipe_id`.
- Partial index `idx_recipe_versions_recipe_published_created_at` supports listing published versions (`published_at IS NOT NULL`) by recipe.
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

## public.audits

Append-only audit trail for domain actions and outcomes.

| Column           | Type        | Constraints                             |
|------------------|-------------|-----------------------------------------|
| `id`             | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid()  |
| `correlation_id` | UUID        | NOT NULL                                |
| `actor_id`       | UUID        | NULLABLE, FK → `users(id)`              |
| `resource_type`  | TEXT        | NOT NULL                                |
| `resource_id`    | UUID        | NULLABLE                                |
| `action`         | TEXT        | NOT NULL                                |
| `created_at`     | TIMESTAMPTZ | NOT NULL, DEFAULT now()                 |
| `result`         | TEXT        | NOT NULL                                |
| `old_state`      | JSONB       | NULLABLE                                |
| `new_state`      | JSONB       | NULLABLE                                |
| `reason`         | TEXT        | NULLABLE                                |
| `context`        | JSONB       | NULLABLE                                |

Notes:
- Indexes exist on `actor_id`, `resource_type/resource_id`, `created_at`, and `correlation_id`.

---

## public.processed_events

Idempotency ledger for event-consumer handlers.

| Column           | Type        | Constraints                         |
|------------------|-------------|-------------------------------------|
| `event_id`       | UUID        | NOT NULL                            |
| `consumer_group` | TEXT        | NOT NULL                            |
| `handler_name`   | TEXT        | NOT NULL                            |
| `processed_at`   | TIMESTAMPTZ | NOT NULL, DEFAULT now()             |

Notes:
- Composite primary key: (`event_id`, `consumer_group`, `handler_name`).

---

## reference.ingredients

Read-mostly reference table of known ingredients. Populated via the ingredient sync command.

| Column                 | Type    | Constraints      |
|------------------------|---------|------------------|
| `id`                   | UUID    | PRIMARY KEY      |
| `name`                 | TEXT    | NOT NULL         |
| `preferred_unit`       | INTEGER | NOT NULL         |
| `file_key`             | TEXT    | NOT NULL, UNIQUE |
| `counter`              | TEXT    | NULLABLE         |
| `plural`               | TEXT    | NULLABLE         |
| `counter_plural`       | TEXT    | NULLABLE         |
| `animal_product_level` | INTEGER | NOT NULL, DEFAULT 0 |
| `contains_gluten`      | BOOLEAN | NOT NULL, DEFAULT FALSE |
| `processing_level`     | INTEGER | NOT NULL, DEFAULT 1 |
| `taxonomy_parent_id`   | UUID    | NULLABLE, FK -> `reference.ingredients(id)` ON DELETE SET NULL |
| `is_searchable`        | BOOLEAN | NOT NULL, DEFAULT TRUE |

Notes:
- `file_key` is a stable identifier used to match rows during upsert syncs from `reference/ingredients.yaml`.
- `preferred_unit` is the unit that the service enforces when creating ingredient usages.
- `taxonomy_parent_id` models taxonomy relationships (specificity), not derivation/component lineage.
- Constraint `ingredients_processing_level_check` enforces `processing_level IN (1, 2, 3)`.
- Constraint `ingredients_taxonomy_parent_not_self` prevents self-parenting.
- Index `idx_ingredients_taxonomy_parent_id` supports parent->children lookups.

---

## Key Relationships

```
users
  ├── refresh_tokens (user_id)
  │     └── refresh_tokens (replaced_by_token_id)
  ├── uploads (owner_user_id)
  ├── audits (actor_id)
  └── recipe_containers (user_id)
    └── recipe_versions (recipe_id)
      ├── recipe_sources (recipe_version_id)    [0..1 per version]
      └── ingredient_usages (recipe_version_id) [0..N per version]

uploads
  └── claimed by domain rows via (linked_entity_type, linked_entity_id)

reference.ingredients
  ├── reference.ingredients (taxonomy_parent_id)
  └── ingredient_usages (ingredient_id)
```

## Migration History Summary

| Migration | Description                                       |
|-----------|---------------------------------------------------|
| 0001      | create recipes                                    |
| 0002      | create users                                      |
| 0003      | create ingredients                                |
| 0004      | create ingredient usages                          |
| 0005      | alter recipes add metadata                        |
| 0006      | create recipe sources                             |
| 0007      | alter recipes add deleted on                      |
| 0008      | alter ingredients add counter                     |
| 0009      | alter ingredients add plurals                     |
| 0010      | alter users add username                          |
| 0011      | alter recipes rename recipe versions              |
| 0012      | alter table rename recipe references              |
| 0013      | alter table recipe versions drop user and deleted |
| 0014      | create recipe containers                          |
| 0015      | alter recipe versions add created at              |
| 0016      | alter recipe versions add version number          |
| 0017      | create refresh tokens                             |
| 0018      | alter recipes add created at id index             |
| 0019      | alter recipe versions add search indexes          |
| 0020      | create uploads                                    |
| 0021      | alter recipe versions add img src                 |
| 0022      | alter uploads add deleted at                      |
| 0023      | alter recipe versions add description             |
| 0024      | alter ingredients add animal product level        |
| 0025      | alter recipe versions add animal product level    |
| 0026      | alter ingredients add contains gluten             |
| 0027      | alter recipe versions add contains gluten         |
| 0028      | create audits                                     |
| 0029      | create processed events                           |
| 0030      | alter ingredients add taxonomy                    |
| 0031      | alter recipes add published at                    |

