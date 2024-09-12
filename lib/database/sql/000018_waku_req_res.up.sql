ALTER TABLE IF NOT EXISTS wakuRequestResponse (
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