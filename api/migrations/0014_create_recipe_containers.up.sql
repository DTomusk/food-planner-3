CREATE TABLE recipe_containers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_version_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_on TIMESTAMPTZ,

    CONSTRAINT fk_current_version
        FOREIGN KEY (current_version_id)
        REFERENCES recipe_versions(id)
        DEFERRABLE INITIALLY DEFERRED
);

ALTER TABLE recipe_versions
    ADD COLUMN recipe_id UUID NOT NULL REFERENCES recipe_containers(id) ON DELETE CASCADE;