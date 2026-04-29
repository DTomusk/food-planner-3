ALTER TABLE recipe_containers
DROP CONSTRAINT fk_draft_version;

ALTER TABLE recipe_containers
DROP COLUMN draft_version_id;