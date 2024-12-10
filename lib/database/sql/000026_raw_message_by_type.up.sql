CREATE TABLE IF NOT EXISTS rawMessageByType (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    size INTEGER NOT NULL,
    messageType TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT rawMessageByType_unique UNIQUE (recordId, timestamp)
);

ALTER TABLE rawMessageByType ADD CONSTRAINT fk_rawMessageByType_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);