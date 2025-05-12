-- 004_routes.sql  ---------------------------------------------------------
CREATE TABLE IF NOT EXISTS routes (
                                      id           bigserial PRIMARY KEY,
                                      from_code    text NOT NULL,
                                      to_code      text NOT NULL,
                                      date         date NOT NULL,
                                      json         jsonb NOT NULL,
                                      cached_at    timestamptz DEFAULT now()
    );
CREATE INDEX IF NOT EXISTS routes_key_idx
    ON routes (from_code, to_code, date);

-- 005_trains.sql  ---------------------------------------------------------
CREATE TABLE IF NOT EXISTS trains (
                                      uid          text PRIMARY KEY,
                                      last_status  jsonb,
                                      last_update  timestamptz
);

-- 006_subscriptions.sql  --------------------------------------------------
CREATE TABLE IF NOT EXISTS subscriptions (
                                             id           bigserial PRIMARY KEY,
                                             device_token text NOT NULL,
                                             train_uid    text NOT NULL,
                                             created_at   timestamptz DEFAULT now(),
    UNIQUE (device_token, train_uid)
    );

-- 007_events.sql  ---------------------------------------------------------
CREATE TABLE IF NOT EXISTS events (
                                      id           bigserial PRIMARY KEY,
                                      ts           timestamptz DEFAULT now(),
    device_id    text,
    event_type   text NOT NULL,
    payload      jsonb
    );

-- 008_extension_pgnotify.sql   --------------------------------
CREATE EXTENSION IF NOT EXISTS pg_notify;  -- для push через LISTEN/NOTIFY
