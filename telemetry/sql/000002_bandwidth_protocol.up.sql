CREATE TABLE IF NOT EXISTS protocolStatsRate (
    id SERIAL PRIMARY KEY,
    hostId VARCHAR(255) NOT NULL,
    protocolName VARCHAR(255) NOT NULL,
    rateIn DOUBLE PRECISION NOT NULL,
    rateOut DOUBLE PRECISION NOT NULL,
    createdAt INTEGER NOT NULL,
    constraint protocolStatsRate_unique unique(hostId, protocolName, createdAt)
);

CREATE TABLE IF NOT EXISTS protocolStatsTotals (
    id SERIAL PRIMARY KEY,
    hostId VARCHAR(255) NOT NULL,
    protocolName VARCHAR(255) NOT NULL,
    totalIn INTEGER NOT NULL,
    totalOut INTEGER NOT NULL,
    createdAt DATE,
    constraint protocolStatsTotals_unique unique(hostId, protocolName, createdAt)
);

ALTER TABLE receivedMessages ADD COLUMN messageSize INTEGER;
