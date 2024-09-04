ALTER TABLE receivedMessages ADD COLUMN statusVersion VARCHAR(31);
ALTER TABLE receivedEnvelopes ADD COLUMN statusVersion VARCHAR(31);
ALTER TABLE sentEnvelopes ADD COLUMN statusVersion VARCHAR(31);