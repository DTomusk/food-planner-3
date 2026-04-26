DROP INDEX IF EXISTS idx_recipe_versions_recipe_published_created_at;
DROP INDEX IF EXISTS ux_recipe_versions_one_draft_per_recipe;

ALTER TABLE recipe_versions
DROP COLUMN IF EXISTS published_at;

DROP INDEX IF EXISTS idx_recipe_containers_public_created_at_id;

ALTER TABLE recipe_containers
DROP COLUMN IF EXISTS published_at;