CREATE TABLE IF NOT EXISTS stations (
                                        code            TEXT PRIMARY KEY,
                                        title           TEXT NOT NULL,
                                        station_type    TEXT,
                                        transport_type  TEXT,
                                        latitude        DOUBLE PRECISION,
                                        longitude       DOUBLE PRECISION,
                                        settlement_code TEXT
);
