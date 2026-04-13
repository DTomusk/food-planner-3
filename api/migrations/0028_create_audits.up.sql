CREATE TABLE audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    correlation_id UUID NOT NULL,
    actor_id UUID,
    resource_type TEXT NOT NULL,
    resource_id UUID,
    action TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    result TEXT NOT NULL,
    old_state JSONB,
    new_state JSONB,
    reason TEXT,
    context JSONB,

    FOREIGN KEY (actor_id) REFERENCES users(id)
);

CREATE INDEX idx_audits_actor_id ON audits(actor_id);
CREATE INDEX idx_audits_resource ON audits(resource_type, resource_id);
CREATE INDEX idx_audits_created_at ON audits(created_at);
CREATE INDEX idx_audits_correlation_id ON audits(correlation_id);