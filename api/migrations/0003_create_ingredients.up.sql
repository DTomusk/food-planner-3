CREATE SCHEMA IF NOT EXISTS reference;

CREATE TABLE reference.ingredients (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    preferred_unit INTEGER NOT NULL,
    file_key TEXT NOT NULL UNIQUE
);