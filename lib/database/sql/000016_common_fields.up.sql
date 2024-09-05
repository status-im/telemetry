-- Create telemetryRecord table
CREATE TABLE IF NOT EXISTS telemetryRecord (
    id SERIAL PRIMARY KEY,
    nodeName VARCHAR(255),
    peerId VARCHAR(255),
    statusVersion VARCHAR(255),
    deviceType VARCHAR(255),
    createdAt INTEGER DEFAULT EXTRACT(EPOCH FROM NOW())::INTEGER
);

-- Function to add recordId column
CREATE OR REPLACE FUNCTION add_telemetryrecord_column(table_name TEXT) RETURNS VOID AS $$
BEGIN
    EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS recordId INTEGER', table_name);
    RAISE NOTICE 'Added recordId column to table: %', table_name;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error adding recordId column to table %: %', table_name, SQLERRM;
END;
$$ LANGUAGE plpgsql;

-- Migration function for existing records
CREATE OR REPLACE FUNCTION migrate_existing_records(table_name TEXT) RETURNS VOID AS $$
BEGIN
    EXECUTE format('
        WITH inserted_telemetry_record AS (
            INSERT INTO telemetryRecord (nodeName, peerId, statusVersion, deviceType, createdAt)
            SELECT DISTINCT nodeName, peerId, statusVersion, deviceType, createdAt
            FROM %I
            RETURNING id, nodeName, peerId, statusVersion, deviceType, createdAt
        )
        UPDATE %I t
        SET recordId = icf.id
        FROM inserted_telemetry_record icf
        WHERE t.nodeName = icf.nodeName
          AND t.peerId = icf.peerId
          AND t.statusVersion = icf.statusVersion
          AND t.deviceType = icf.deviceType
          AND t.createdAt = icf.createdAt;
    ', table_name, table_name);

    RAISE NOTICE 'Completed migration for table: %', table_name;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error migrating data for table %: %', table_name, SQLERRM;
END;
$$ LANGUAGE plpgsql;

-- Function to modify tables
CREATE OR REPLACE FUNCTION modify_table(table_name TEXT) RETURNS VOID AS $$
BEGIN
    -- Create the foreign key constraint
    BEGIN
        EXECUTE format('ALTER TABLE %I ADD CONSTRAINT fk_%I_telemetryRecord
            FOREIGN KEY (recordId) REFERENCES telemetryRecord(id)', table_name, table_name);
    EXCEPTION
        WHEN duplicate_object THEN
            RAISE NOTICE 'Foreign key constraint already exists on %', table_name;
    END;

    -- Remove columns that are now in telemetryRecord
    EXECUTE format('ALTER TABLE %I
        DROP COLUMN IF EXISTS createdAt,
        DROP COLUMN IF EXISTS nodeName,
        DROP COLUMN IF EXISTS statusVersion,
        DROP COLUMN IF EXISTS peerId,
        DROP COLUMN IF EXISTS deviceType', table_name);

    -- Drop the unique constraint
    EXECUTE format('ALTER TABLE %I DROP CONSTRAINT IF EXISTS %I_unique', table_name, table_name);

    -- Make the new column NOT NULL
    EXECUTE format('ALTER TABLE %I ALTER COLUMN recordId SET NOT NULL', table_name);

    RAISE NOTICE 'Completed modifications for table: %', table_name;
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error modifying table %: %', table_name, SQLERRM;
END;
$$ LANGUAGE plpgsql;

-- Add recordId column to each table
SELECT add_telemetryrecord_column('peercount');
SELECT add_telemetryrecord_column('receivedmessages');
SELECT add_telemetryrecord_column('receivedenvelopes');
SELECT add_telemetryrecord_column('sentenvelopes');
SELECT add_telemetryrecord_column('errorsendingenvelope');
SELECT add_telemetryrecord_column('peerconnfailure');

-- Apply migration to each table
SELECT migrate_existing_records('peercount');
SELECT migrate_existing_records('receivedmessages');
SELECT migrate_existing_records('receivedenvelopes');
SELECT migrate_existing_records('sentenvelopes');
SELECT migrate_existing_records('errorsendingenvelope');
SELECT migrate_existing_records('peerconnfailure');

-- Apply modifications to each table
SELECT modify_table('peercount');
SELECT modify_table('receivedmessages');
SELECT modify_table('receivedenvelopes');
SELECT modify_table('sentenvelopes');
SELECT modify_table('errorsendingenvelope');
SELECT modify_table('peerconnfailure');

-- Drop the functions after use
DROP FUNCTION add_telemetryrecord_column;
DROP FUNCTION migrate_existing_records;
DROP FUNCTION modify_table;

-- Recreate unique constraints
ALTER TABLE receivedMessages 
ADD CONSTRAINT receivedMessages_unique 
UNIQUE (
    recordId,
    chatId, 
    messageHash, 
    receiverKeyUID 
);

ALTER TABLE receivedEnvelopes 
ADD CONSTRAINT receivedEnvelopes_unique 
UNIQUE (
    recordId,
    sentAt,
    messageHash, 
    receiverKeyUID
);

ALTER TABLE sentEnvelopes 
ADD CONSTRAINT sentEnvelopes_unique 
UNIQUE (
    recordId,
    sentAt,
    messageHash, 
    senderKeyUID
);

ALTER TABLE errorSendingEnvelope 
ADD CONSTRAINT errorSendingEnvelope_unique 
UNIQUE (
    recordId,
    sentAt,
    messageHash, 
    senderKeyUID,
    error
);

ALTER TABLE peerCount 
ADD CONSTRAINT peerCount_unique 
UNIQUE (
    recordId,
    timestamp,
    nodeKeyUID
);

ALTER TABLE peerConnFailure 
ADD CONSTRAINT peerConnFailure_unique 
UNIQUE (
    recordId,
    timestamp, 
    failedPeerId, 
    failureCount
);