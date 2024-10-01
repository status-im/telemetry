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
		"missingmessages",
		"missingrelevantmessages",
		"messageDeliveryConfirmed",
		"schema_migrations",
	}

	for _, table := range tables {
		_, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Fatalf("an error '%s' was not expected when dropping the table %s", err, table)
		}

	}

	_, err = tx.Exec("DROP INDEX IF EXISTS receivedEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS receivedMessageAggregated_runAt")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS protocolStatsRate_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS protocolStatsTotals_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS peerCount_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS receivedMessages_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS peerConnFailure_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS peerCountByShard_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS peerCountByOrigin_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS missingMessages_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS messageCheckSuccess_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS messageCheckFailure_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS missingRelevantMessages_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS messageDeliveryConfirmed_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = tx.Exec("DROP INDEX IF EXISTS dialFailure_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
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
