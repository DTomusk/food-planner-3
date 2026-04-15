CREATE TABLE processed_events (
    event_id UUID NOT NULL, 
    consumer_group TEXT NOT NULL,
    handler_name TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_processed_events PRIMARY KEY (event_id, consumer_group, handler_name)
);
