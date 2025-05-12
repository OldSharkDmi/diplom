ALTER TABLE trains ADD COLUMN IF NOT EXISTS id BIGSERIAL PRIMARY KEY;
CREATE TABLE IF NOT EXISTS train_runs (
                                          id            BIGSERIAL PRIMARY KEY,
                                          train_id      BIGINT      NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    run_date      DATE        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (train_id, run_date)
    );

CREATE TABLE IF NOT EXISTS train_statuses (
                                              id            BIGSERIAL PRIMARY KEY,
                                              train_run_id  BIGINT      NOT NULL REFERENCES train_runs(id) ON DELETE CASCADE,
    status        JSONB       NOT NULL,
    fetched_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (train_run_id)
    );

CREATE TABLE IF NOT EXISTS occupancy_predictions (
                                                     id              BIGSERIAL PRIMARY KEY,
                                                     train_run_id    BIGINT      NOT NULL REFERENCES train_runs(id) ON DELETE CASCADE,
    car_number      INT         NOT NULL,
    level           TEXT        NOT NULL CHECK (level IN ('low', 'medium', 'high')),
    predicted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (train_run_id, car_number)
    );