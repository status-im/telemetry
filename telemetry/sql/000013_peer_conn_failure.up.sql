CREATE TABLE IF NOT EXISTS peerConnFailure (
    id SERIAL PRIMARY KEY,
    createdAt INTEGER NOT NULL,
    peerId VARCHAR(255) NOT NULL,
    nodeName VARCHAR(255) NOT NULL,
    nodeKeyUid VARCHAR(255) NOT NULL,
    timestamp INTEGER NOT NULL,
    statusVersion VARCHAR(31),
    failureCount INTEGER NOT NULL,
    failedPeerId VARCHAR(255) NOT NULL,
    CONSTRAINT peerConnFailure_unique unique(timestamp, peerId, failedPeerId, failureCount)
);
