DO $$ BEGIN
    CREATE TYPE publish_method AS ENUM ('LightPush', 'Relay');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS sentEnvelopes (
    id SERIAL PRIMARY KEY,
    messageHash VARCHAR(255) NOT NULL,
    sentAt INTEGER NOT NULL,
    createdAt INTEGER NOT NULL,
    topic VARCHAR(255) NOT NULL,
    pubSubTopic VARCHAR(255) NOT NULL,
    senderKeyUID VARCHAR(255) NOT NULL,
    nodeName VARCHAR(255) NOT NULL,
    publishMethod publish_method,
    CONSTRAINT sentEnvelopes_unique unique(sentAt, messageHash, senderKeyUID, nodeName)
);
