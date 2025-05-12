-- 009_trains_notify.sql -------------  PUSH через LISTEN/NOTIFY
CREATE OR REPLACE FUNCTION notify_train_update() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    PERFORM pg_notify('train_updates', NEW.uid);
RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_train_notify ON trains;
CREATE TRIGGER trg_train_notify
    AFTER UPDATE OF last_status ON trains
    FOR EACH ROW
    WHEN (OLD.last_status IS DISTINCT FROM NEW.last_status)
EXECUTE PROCEDURE notify_train_update();

-- 010_events_index.sql --------------  быстрый поиск по типу
CREATE INDEX IF NOT EXISTS events_type_idx ON events (event_type);

-- 011_events_gin.sql ----------------  full-text внутри payload
CREATE INDEX IF NOT EXISTS events_payload_gin
    ON events USING gin (payload jsonb_path_ops);

-- 012_trains_uid_idx.sql ------------  (если ещё нет)
CREATE UNIQUE INDEX IF NOT EXISTS trains_uid_idx ON trains(uid);
