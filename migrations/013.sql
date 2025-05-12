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
    status        VARCHAR(50) NOT NULL,
    received_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw           JSONB,
    UNIQUE (train_run_id)
    );

CREATE TYPE occupancy_level AS ENUM ('low', 'medium', 'high');

CREATE TABLE IF NOT EXISTS occupancy_predictions (
                                                     id             BIGSERIAL PRIMARY KEY,
                                                     train_run_id   BIGINT            NOT NULL REFERENCES train_runs(id) ON DELETE CASCADE,
    car_number     SMALLINT          NOT NULL,
    level          occupancy_level   NOT NULL,
    predicted_at   TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    UNIQUE (train_run_id, car_number)
    );