ALTER TABLE peercount ADD COLUMN deviceType VARCHAR(255);
ALTER TABLE receivedMessages ADD COLUMN deviceType VARCHAR(255);
ALTER TABLE receivedEnvelopes ADD COLUMN deviceType VARCHAR(255);
ALTER TABLE sentEnvelopes ADD COLUMN deviceType VARCHAR(255);
ALTER TABLE errorSendingEnvelope ADD COLUMN deviceType VARCHAR(255);
ALTER TABLE peerConnFailure ADD COLUMN deviceType VARCHAR(255);