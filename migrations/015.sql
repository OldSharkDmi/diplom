DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'routes_from_to_date_unique'
    ) THEN
ALTER TABLE routes
    ADD CONSTRAINT routes_from_to_date_unique
        UNIQUE (from_code, to_code, date);
END IF;
END$$;
