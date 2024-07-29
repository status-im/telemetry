CREATE TABLE IF NOT EXISTS wakuPushError (
	id SERIAL PRIMARY KEY,
	peerId VARCHAR(255) NOT NULL,
	errorMessage VARCHAR(255) NOT NULL,
	peerIdRemote VARCHAR(255),
	contentTopic VARCHAR(255),
	pubsubTopic VARCHAR(255),
	timestamp INTEGER NOT NULL,
	createdAt INTEGER NOT NULL,

    CONSTRAINT wakuPushError_unique unique(peerId, errorMessage, timestamp, peerIdRemote)
);