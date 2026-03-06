ALTER TABLE recipe_versions ADD COLUMN user_id UUID REFERENCES users(id);
ALTER TABLE recipe_versions ADD COLUMN deleted_on TIMESTAMPTZ;