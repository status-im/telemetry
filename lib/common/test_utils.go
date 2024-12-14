package common

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/status-im/telemetry/lib/database"
)

func NewMock() *sql.DB {
	db, err := sql.Open("postgres", "postgres://telemetry:newPassword@127.0.0.1:5432/telemetrydb?sslmode=disable")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	err = database.CreateTables(db)
	if err != nil {
		log.Fatalf("an error '%s' was not expected when migrating the db", err)
	}

	return db
}

func DropTables(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	tables := []string{
		"receivedMessages",
		"receivedMessageAggregated",
		"receivedEnvelopes",
		"sentEnvelopes",
		"protocolStatsRate",
		"protocolStatsTotals",
		"peercount",
		"peerconnfailure",
		"errorsendingenvelope",
		"peerCountByShard",
		"peerCountByOrigin",
		"messageCheckSuccess",
		"messageCheckFailure",
		"dialFailure",
		"missedMessages",
		"missedrelevantmessages",
		"messageDeliveryConfirmed",
		"sentMessageTotal",
		"rawMessageByType",
		"schema_migrations",
	}

	for _, table := range tables {
		_, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Fatalf("an error '%s' was not expected when dropping the table %s", err, table)
		}

	}

	indexes := []string{
		"receivedEnvelopes",
		"receivedMessageAggregated_runAt",
		"protocolStatsTotals_idx1",
		"peerCount_unique",
		"receivedMessages_unique",
		"peerConnFailure_unique",
		"peerCountByShard_unique",
		"peerCountByOrigin_unique",
		"missedMessages_unique",
		"messageCheckSuccess_unique",
		"messageCheckFailure_unique",
		"missedRelevantMessages_unique",
		"messageDeliveryConfirmed_unique",
		"dialFailure_unique",
		"sentMessageTotal_unique",
		"rawMessageByType_unique",
	}

	for _, index := range indexes {
		_, err := tx.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %s", index))
		if err != nil {
			log.Fatalf("an error '%s' was not expected when dropping the index", err)
		}
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS schema_migrations")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS telemetryRecord")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit the TX: %s", err)
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("failed to close db: %s", err)
	}
}
