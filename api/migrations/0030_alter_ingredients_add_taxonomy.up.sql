ALTER TABLE reference.ingredients
ADD COLUMN IF NOT EXISTS processing_level INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS taxonomy_parent_id UUID REFERENCES reference.ingredients(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS is_searchable BOOLEAN NOT NULL DEFAULT TRUE,
ADD CONSTRAINT ingredients_processing_level_check CHECK (processing_level IN (1, 2, 3)),
ADD CONSTRAINT ingredients_taxonomy_parent_not_self CHECK (taxonomy_parent_id IS NULL OR taxonomy_parent_id <> id);

CREATE INDEX idx_ingredients_taxonomy_parent_id
ON reference.ingredients(taxonomy_parent_id);