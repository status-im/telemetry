CREATE TABLE IF NOT EXISTS wakuPushFilter (
	id SERIAL PRIMARY KEY,
	walletAddress VARCHAR(255),
	peerIdSender VARCHAR(255) NOT NULL,
	peerIdReporter VARCHAR(255) NOT NULL,
	sequenceHash VARCHAR(255) NOT NULL,
	sequenceTotal VARCHAR(255) NOT NULL,
	sequenceIndex VARCHAR(255) NOT NULL,
	contentTopic VARCHAR(255) NOT NULL,
	pubsubTopic VARCHAR(255) NOT NULL,
	timestamp INTEGER NOT NULL,
	createdAt INTEGER NOT NULL,

    CONSTRAINT wakuPushFilter_unique unique(peerIdSender, peerIdReporter, sequenceHash, sequenceIndex)
);