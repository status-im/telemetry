CREATE TABLE IF NOT EXISTS peerCount (
    id SERIAL PRIMARY KEY,
    createdAt INTEGER NOT NULL,
    peerCount INTEGER NOT NULL,
    nodeName VARCHAR(255) NOT NULL,
    nodeKeyUid VARCHAR(255) NOT NULL,
    timestamp INTEGER NOT NULL,
    statusVersion VARCHAR(31),
    CONSTRAINT peerCount_unique unique(timestamp, nodeName, nodeKeyUid, statusVersion)
);
