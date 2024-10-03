CREATE TABLE IF NOT EXISTS missedMessages (
    id SERIAL PRIMARY KEY,
    messageHash TEXT NOT NULL,
    sentAt INTEGER NOT NULL,
    contentTopic TEXT NOT NULL,
    pubsubTopic TEXT NOT NULL,  
    recordId INTEGER NOT NULL
);

ALTER TABLE missedMessages ADD CONSTRAINT fk_missedMessages_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

ALTER TABLE missedMessages 
ADD CONSTRAINT missedMessages_unique 
UNIQUE (
    recordId,
    messageHash,
    contentTopic
);

CREATE TABLE IF NOT EXISTS missedRelevantMessages (
    id SERIAL PRIMARY KEY,
    messageHash TEXT NOT NULL,
    sentAt INTEGER NOT NULL,
    contentTopic TEXT NOT NULL,
    pubsubTopic TEXT NOT NULL,  
    recordId INTEGER NOT NULL
);

ALTER TABLE missedRelevantMessages ADD CONSTRAINT fk_missedRelevantMessages_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

ALTER TABLE missedRelevantMessages 
ADD CONSTRAINT missedRelevantMessages_unique 
UNIQUE (
    recordId,
    messageHash,
    contentTopic
);