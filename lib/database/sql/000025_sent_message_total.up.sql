CREATE TABLE IF NOT EXISTS sentMessageTotal (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    size INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT sentMessageTotal_unique UNIQUE (recordId, timestamp)
);

ALTER TABLE sentMessageTotal ADD CONSTRAINT fk_sentMessageTotal_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);