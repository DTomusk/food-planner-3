CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

ALTER TABLE recipes ADD COLUMN user_id UUID REFERENCES users(id);