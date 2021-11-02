package telemetry

import (
	"database/sql"
	"log"

	// TODO Replace with real db driver
	_ "github.com/mattn/go-sqlite3" // Blank import to register the sqlite3 driver
)

func OpenDb(username string, password string, name string) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
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
		id INTEGER PRIMARY KEY,
		chatId VARCHAR(255) NOT NULL,
		messageHash VARCHAR(255) NOT NULL,
		receiverKeyUID VARCHAR(255) NOT NULL,
		sentAt INTEGER NOT NULL,
		topic VARCHAR(255) NOT NULL,
		createdAt INTEGER NOT NULL
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
