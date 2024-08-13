TRUNCATE TABLE receivedMessages;
TRUNCATE TABLE protocolStatsRate;
TRUNCATE TABLE protocolStatsTotals;
TRUNCATE TABLE receivedMessageAggregated;

CREATE INDEX receivedMessages_createdAt ON receivedMessages(createdAt);
CREATE INDEX receivedMessageAggregated_runAt ON receivedMessageAggregated(runAt);

CREATE INDEX protocolStatsRate_idx1 ON protocolStatsRate(protocolName, createdAt);
CREATE INDEX protocolStatsTotals_idx1 ON protocolStatsTotals(protocolName, createdAt);

ALTER TABLE protocolStatsRate RENAME COLUMN hostID TO peerID;
ALTER TABLE protocolStatsTotals RENAME COLUMN hostID TO peerID;
