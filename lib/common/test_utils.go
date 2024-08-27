package common

import (
	"database/sql"
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
	_, err := db.Exec("DROP TABLE IF EXISTS receivedMessages")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS receivedMessageAggregated")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS receivedEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS sentEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS protocolStatsRate")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS protocolStatsTotals")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS peercount")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS peerconnfailure")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS errorsendingenvelope")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS schema_migrations")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the table", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS receivedEnvelopes")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS receivedMessageAggregated_runAt")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS protocolStatsRate_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS protocolStatsTotals_idx1")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS peerCount_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS receivedMessages_unique")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	_, err = db.Exec("DROP INDEX IF EXISTS peerConnFailure_uniqu")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when dropping the index", err)
	}

	db.Close()
}
