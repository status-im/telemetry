CREATE TABLE IF NOT EXISTS wakuGeneric (
	id SERIAL PRIMARY KEY,
	peerId VARCHAR(255) NOT NULL,
	metricType VARCHAR(255) NOT NULL,
	
	contentTopic VARCHAR(255),
	pubsubTopic VARCHAR(255),
	genericData VARCHAR(255),
	errorMessage VARCHAR(255),

	timestamp INTEGER NOT NULL,
	createdAt INTEGER NOT NULL,

    CONSTRAINT wakuGeneric_unique unique(timestamp, metricType, peerId)
);