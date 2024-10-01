CREATE TABLE IF NOT EXISTS dial_error_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dialFailure (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    errorType INTEGER NOT NULL,
    errorMsg TEXT,
    protocols TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT dialFailure_unique UNIQUE (recordId, errorType, protocols, timestamp),
    CONSTRAINT fk_dialFailure_errorType FOREIGN KEY (errorType) REFERENCES dial_error_types(id)
);

ALTER TABLE dialFailure ADD CONSTRAINT fk_dialFailure_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

INSERT INTO dial_error_types (id, name)
SELECT v.id, v.name
FROM (VALUES
    (0, 'Unknown'),
    (1, 'I/O Timeout'),
    (2, 'Connection Refused'),
    (3, 'Relay Circuit Failed'),
    (4, 'Relay No Reservation'),
    (5, 'Security Negotiation Failed'),
    (6, 'Concurrent Dial Succeeded'),
    (7, 'Concurrent Dial Failed')
) AS v(id, name)
WHERE NOT EXISTS (
    SELECT 1 FROM dial_error_types WHERE id = v.id
);