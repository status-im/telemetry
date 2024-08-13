CREATE TABLE IF NOT EXISTS receivedEnvelopes (
    id SERIAL PRIMARY KEY,
    messageHash VARCHAR(255) NOT NULL,
    sentAt INTEGER NOT NULL,
    createdAt INTEGER NOT NULL,
    topic VARCHAR(255) NOT NULL,
    pubSubTopic VARCHAR(255) NOT NULL,
    receiverKeyUID VARCHAR(255) NOT NULL,
    nodeName VARCHAR(255) NOT NULL,
    processingError VARCHAR(255) NOT NULL,
    CONSTRAINT receivedEnvelopes_unique unique(sentAt, messageHash, receiverKeyUID, nodeName)
);

ALTER TABLE receivedMessages ADD COLUMN pubSubTopic VARCHAR(255);
