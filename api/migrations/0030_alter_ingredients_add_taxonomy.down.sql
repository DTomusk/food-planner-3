DROP INDEX IF EXISTS idx_ingredients_taxonomy_parent_id;

ALTER TABLE reference.ingredients
DROP CONSTRAINT IF EXISTS ingredients_processing_level_check,
DROP CONSTRAINT IF EXISTS ingredients_taxonomy_parent_not_self,
DROP COLUMN IF EXISTS processing_level,
DROP COLUMN IF EXISTS taxonomy_parent_id,
DROP COLUMN IF EXISTS is_searchable;