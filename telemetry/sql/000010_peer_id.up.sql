ALTER TABLE peercount ADD COLUMN peerId VARCHAR(255);
ALTER TABLE receivedMessages ADD COLUMN peerId VARCHAR(255);
ALTER TABLE receivedEnvelopes ADD COLUMN peerId VARCHAR(255);
ALTER TABLE sentEnvelopes ADD COLUMN peerId VARCHAR(255);
ALTER TABLE errorSendingEnvelope ADD COLUMN peerId VARCHAR(255);