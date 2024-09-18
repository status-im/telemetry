CREATE TABLE IF NOT EXISTS wakuRequestResponse (
	id SERIAL PRIMARY KEY,
	protocol VARCHAR(50) NOT NULL,
	ephemeral BOOLEAN NOT NULL,
	timestamp INTEGER NOT NULL,
	seenTimestamp INTEGER NOT NULL,
	createdAt INTEGER NOT NULL,
	contentTopic VARCHAR(255) NOT NULL,
	pubsubTopic VARCHAR(255) NOT NULL,
	peerId VARCHAR(255) NOT NULL,
	messageHash VARCHAR(255) NOT NULL,
	errorMessage TEXT,
	extraData TEXT,

	CONSTRAINT messages_unique UNIQUE (peerId, messageHash)
);

CREATE TABLE IF NOT EXISTS messageCheckSuccess (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    messageHash TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT messageCheckSuccess_unique UNIQUE (recordId, messageHash, timestamp)
);

ALTER TABLE messageCheckSuccess ADD CONSTRAINT fk_messageCheckSuccess_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);
            
CREATE TABLE IF NOT EXISTS messageCheckFailure (
    id SERIAL PRIMARY KEY,
    recordId INTEGER NOT NULL,
    messageHash TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    CONSTRAINT messageCheckFailure_unique UNIQUE (recordId, messageHash, timestamp)
);

ALTER TABLE messageCheckFailure ADD CONSTRAINT fk_messageCheckFailure_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id);