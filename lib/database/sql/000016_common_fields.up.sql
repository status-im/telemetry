CREATE TABLE IF NOT EXISTS commonFields (
    id SERIAL PRIMARY KEY,
    nodeName VARCHAR(255) NOT NULL,
    peerId VARCHAR(255) NOT NULL,
    statusVersion VARCHAR(255) NOT NULL,
    deviceType VARCHAR(255) NOT NULL,
    createdAt INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())::INTEGER
);

-- Function to modify tables with added logging
CREATE OR REPLACE FUNCTION modify_table(table_name TEXT) RETURNS VOID AS $$
BEGIN

    -- Add the new column for the foreign key
    EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS commonFieldsId INTEGER', table_name);

    -- Create the foreign key constraint
    BEGIN
        EXECUTE format('ALTER TABLE %I ADD CONSTRAINT fk_%I_commonFields
            FOREIGN KEY (commonFieldsId) REFERENCES commonFields(id)', table_name, table_name);
    EXCEPTION
        WHEN duplicate_object THEN
            RAISE NOTICE 'Foreign key constraint already exists on %', table_name;
    END;

    -- Remove columns that are now in commonFields
    EXECUTE format('ALTER TABLE %I
        DROP COLUMN IF EXISTS createdAt,
        DROP COLUMN IF EXISTS nodeName,
        DROP COLUMN IF EXISTS statusVersion,
        DROP COLUMN IF EXISTS peerId,
        DROP COLUMN IF EXISTS deviceType', table_name);

    -- Drop the unique constraint
    EXECUTE format('ALTER TABLE %I DROP CONSTRAINT IF EXISTS %I_unique', table_name, table_name);

    -- Make the new column NOT NULL
    EXECUTE format('ALTER TABLE %I ALTER COLUMN commonFieldsId SET NOT NULL', table_name);

    -- Add foreign key constraint
    EXECUTE format('
        ALTER TABLE %I
        ADD CONSTRAINT %I_commonfields_fk
        FOREIGN KEY (commonFieldsId) REFERENCES commonFields(id)
    ', table_name, table_name);

    RAISE NOTICE 'Completed modifications for table: %', table_name;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error modifying table %: %', table_name, SQLERRM;
END;
$$ LANGUAGE plpgsql;

-- Apply modifications to each table
SELECT modify_table('peercount');
SELECT modify_table('receivedmessages');
SELECT modify_table('receivedenvelopes');
SELECT modify_table('sentenvelopes');
SELECT modify_table('errorsendingenvelope');
SELECT modify_table('peerconnfailure');

-- Drop the function after use
DROP FUNCTION modify_table;

-- Recreate unique constraints
ALTER TABLE receivedMessages 
ADD CONSTRAINT receivedMessages_unique 
UNIQUE (
    commonFieldsId,
    chatId, 
    messageHash, 
    receiverKeyUID 
);

ALTER TABLE receivedEnvelopes 
ADD CONSTRAINT receivedEnvelopes_unique 
UNIQUE (
    commonFieldsId,
    sentAt,
    messageHash, 
    receiverKeyUID
);

ALTER TABLE sentEnvelopes 
ADD CONSTRAINT sentEnvelopes_unique 
UNIQUE (
    commonFieldsId,
    sentAt,
    messageHash, 
    senderKeyUID
);

ALTER TABLE errorSendingEnvelope 
ADD CONSTRAINT errorSendingEnvelope_unique 
UNIQUE (
    commonFieldsId,
    sentAt,
    messageHash, 
    senderKeyUID,
    error
);

ALTER TABLE peerCount 
ADD CONSTRAINT peerCount_unique 
UNIQUE (
    commonFieldsId,
    timestamp,
    nodeKeyUid
);

ALTER TABLE peerConnFailure 
ADD CONSTRAINT peerConnFailure_unique 
UNIQUE (
    commonFieldsId,
    timestamp, 
    failedPeerId, 
    failureCount
);