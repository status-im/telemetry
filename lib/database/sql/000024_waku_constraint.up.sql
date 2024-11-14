ALTER TABLE wakuRequestResponse DROP CONSTRAINT messages_unique;

ALTER TABLE wakuRequestResponse ADD CONSTRAINT messages_unique UNIQUE (peerId, messageHash, timestamp);