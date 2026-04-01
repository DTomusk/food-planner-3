CREATE TABLE uploads (
    id UUID PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL UNIQUE,
    file_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    purpose TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    linked_entity_type TEXT,
    linked_entity_id UUID,

    CONSTRAINT valid_purpose CHECK (purpose IN ('recipe-images')),
    CONSTRAINT uploads_link_pair_check
        CHECK (
            (linked_entity_type IS NULL AND linked_entity_id IS NULL) OR
            (linked_entity_type IS NOT NULL AND linked_entity_id IS NOT NULL)
        ),
    CONSTRAINT expires_after_created CHECK (expires_at > created_at),
    CONSTRAINT file_size_non_negative CHECK (file_size_bytes >= 0)
);

CREATE INDEX idx_uploads_owner_user_id ON uploads(owner_user_id);
CREATE INDEX idx_uploads_active_unlinked ON uploads(owner_user_id, purpose, created_at)
WHERE used_at IS NULL;