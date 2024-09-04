CREATE TABLE IF NOT EXISTS storeConfirmationFailed (
    id SERIAL PRIMARY KEY,
    messageHash TEXT NOT NULL,
    recordId INTEGER NOT NULL
);

ALTER TABLE storeConfirmationFailed ADD CONSTRAINT fk_storeConfirmationFailed_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

ALTER TABLE storeConfirmationFailed 
ADD CONSTRAINT storeConfirmationFailed_unique 
UNIQUE (
    recordId,
    messageHash
);