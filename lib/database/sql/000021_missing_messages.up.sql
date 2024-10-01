CREATE TABLE IF NOT EXISTS missingMessages (
    id SERIAL PRIMARY KEY,
    messageHash TEXT NOT NULL,
    sentAt INTEGER NOT NULL,
    contentTopic TEXT NOT NULL,
    pubsubTopic TEXT NOT NULL,  
    recordId INTEGER NOT NULL
);

ALTER TABLE missingMessages ADD CONSTRAINT fk_missingMessages_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

ALTER TABLE missingMessages 
ADD CONSTRAINT missingMessages_unique 
UNIQUE (
    recordId,
    messageHash,
    contentTopic
);

CREATE TABLE IF NOT EXISTS missingRelevantMessages (
    id SERIAL PRIMARY KEY,
    messageHash TEXT NOT NULL,
    sentAt INTEGER NOT NULL,
    contentTopic TEXT NOT NULL,
    pubsubTopic TEXT NOT NULL,  
    recordId INTEGER NOT NULL
);

ALTER TABLE missingRelevantMessages ADD CONSTRAINT fk_missingRelevantMessages_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);

ALTER TABLE missingRelevantMessages 
ADD CONSTRAINT missingRelevantMessages_unique 
UNIQUE (
    recordId,
    messageHash,
    contentTopic
);