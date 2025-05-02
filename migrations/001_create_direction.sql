CREATE TABLE IF NOT EXISTS directions (
                                          id            SERIAL PRIMARY KEY,
                                          name          TEXT NOT NULL,
                                          from_station  TEXT,
                                          to_station    TEXT
);
