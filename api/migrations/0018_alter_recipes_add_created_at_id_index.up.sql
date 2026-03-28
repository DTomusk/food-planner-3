CREATE INDEX idx_recipes_created_at_id 
ON recipe_containers (created_at DESC, id DESC)
WHERE deleted_on IS NULL;