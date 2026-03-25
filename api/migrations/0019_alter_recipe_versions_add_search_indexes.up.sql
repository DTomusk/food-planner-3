-- Trigram indexes for recipe_versions.name to improve search performance
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- full text search index on name
CREATE INDEX idx_recipe_versions_name_fts_gin
ON recipe_versions
USING GIN (to_tsvector('english', coalesce(name, '')));

-- trigram index to support fuzzy search
CREATE INDEX idx_recipe_versions_name_trgm_gin
ON recipe_versions
USING GIN (lower(name) gin_trgm_ops);