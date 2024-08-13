CREATE TABLE IF NOT EXISTS errorSendingEnvelope (
    id SERIAL PRIMARY KEY,
    messageHash VARCHAR(255) NOT NULL,
    sentAt INTEGER NOT NULL,
    createdAt INTEGER NOT NULL,
    topic VARCHAR(255) NOT NULL,
    pubSubTopic VARCHAR(255) NOT NULL,
    senderKeyUID VARCHAR(255) NOT NULL,
    nodeName VARCHAR(255) NOT NULL,
    publishMethod publish_method,
    error TEXT NOT NULL,
    statusVersion VARCHAR(31),
    CONSTRAINT errorSendingEnvelope_unique unique(sentAt, messageHash, senderKeyUID, nodeName, error)
);
