package telemetry

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenDb(dataSourceName string) *sql.DB {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	log.Println("Connected to database")

	if err := createTables(db); err != nil {
		log.Fatalf("unable to create the table: %v", err)
	}
	log.Println("DB initialized")

	return db
}

func createTables(db *sql.DB) error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS receivedMessages (
		id SERIAL PRIMARY KEY,
		chatId VARCHAR(255) NOT NULL,
		messageHash VARCHAR(255) NOT NULL,
		messageId VARCHAR(255) NOT NULL,
		receiverKeyUID VARCHAR(255) NOT NULL,
		nodeName VARCHAR(255) NOT NULL,
		sentAt INTEGER NOT NULL,
		topic VARCHAR(255) NOT NULL,
		createdAt INTEGER NOT NULL,

		constraint receivedMessages_unique unique(chatId, messageHash, receiverKeyUID, nodeName)

	);`
	_, err := db.Exec(sqlStmt)

	if err != nil {
		return err
	}

	sqlStmt = `CREATE TABLE IF NOT EXISTS receivedMessageAggregated (
		id SERIAL PRIMARY KEY,
		durationInSeconds INTEGER NOT NULL,
		chatId VARCHAR(255) NOT NULL,
		value DECIMAL NOT NULL,
		runAt INTEGER NOT NULL
	);`

	_, err = db.Exec(sqlStmt)

	return err
}
