ALTER TABLE recipe_containers
ADD COLUMN published_at timestamptz;

UPDATE recipe_containers rc
SET published_at = rc.created_at;

-- Make fetching published not deleted recipes more efficient
CREATE INDEX idx_recipe_containers_public_created_at_id
ON recipe_containers (created_at DESC, id DESC)
WHERE deleted_on IS NULL AND published_at IS NOT NULL;

ALTER TABLE recipe_versions
ADD COLUMN published_at timestamptz;

UPDATE recipe_versions rv
SET published_at = rv.created_at;

-- do not allow multiple drafts per recipe
CREATE UNIQUE INDEX ux_recipe_versions_one_draft_per_recipe
ON recipe_versions (recipe_id)
WHERE published_at IS NULL;

-- make fetching published versions for a recipe more efficient
CREATE INDEX idx_recipe_versions_recipe_published_created_at
ON recipe_versions (recipe_id, created_at DESC)
WHERE published_at IS NOT NULL;