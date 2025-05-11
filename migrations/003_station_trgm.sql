CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS stations_title_trgm
    ON stations USING GIN (title gin_trgm_ops);
