CREATE TABLE IF NOT EXISTS peerCountByShard (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    count INTEGER NOT NULL,
    shard INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT peerCountByShard_unique UNIQUE (recordId, count, shard, timestamp)
);

ALTER TABLE peerCountByShard ADD CONSTRAINT fk_peerCountByShard_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

CREATE TABLE IF NOT EXISTS peerCountByOrigin (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    count INTEGER NOT NULL,
    origin INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT peerCountByOrigin_unique UNIQUE (recordId, count, origin, timestamp)
);

ALTER TABLE peerCountByOrigin ADD CONSTRAINT fk_peerCountByOrigin_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

CREATE TABLE IF NOT EXISTS origin_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO origin_types (id, name)
SELECT v.id, v.name
FROM (VALUES
    (0, 'Unknown'),
    (1, 'Discv5'),
    (2, 'Static'),
    (3, 'PeerExchange'),
    (4, 'DNSDiscovery'),
    (5, 'Rendezvous'),
    (6, 'PeerManager')
) AS v(id, name)
WHERE NOT EXISTS (
    SELECT 1 FROM origin_types WHERE id = v.id
);