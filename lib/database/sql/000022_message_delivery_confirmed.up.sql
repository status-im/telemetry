CREATE TABLE IF NOT EXISTS messageDeliveryConfirmed (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    messageHash TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT messageDeliveryConfirmed_unique UNIQUE (recordId, messageHash, timestamp)
);

ALTER TABLE messageDeliveryConfirmed ADD CONSTRAINT fk_messageDeliveryConfirmed_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);