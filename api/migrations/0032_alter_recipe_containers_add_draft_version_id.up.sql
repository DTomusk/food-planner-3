ALTER TABLE recipe_containers
ADD COLUMN draft_version_id UUID;

ALTER TABLE recipe_containers
ADD CONSTRAINT fk_draft_version
    FOREIGN KEY (draft_version_id)
    REFERENCES recipe_versions(id)
    DEFERRABLE INITIALLY DEFERRED;