package telemetry

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Migrate applies migrations.
func Migrate(db *sql.DB, driver database.Driver) error {
	return migrateDB(db, bindata.Resource(
		AssetNames(),
		Asset,
	), driver)
}

// Migrate database using provided resources.
func migrateDB(db *sql.DB, resources *bindata.AssetSource, driver database.Driver) error {
	source, err := bindata.WithInstance(resources)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"go-bindata",
		source,
		"telemetrydb",
		driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func OpenDb(dataSourceName string, logger *zap.Logger) *sql.DB {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	logger.Info("Connected to database")

	if err := createTables(db); err != nil {
		log.Fatalf("unable to create the table: %v", err)
	}
	logger.Info("DB initialized")

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
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: postgres.DefaultMigrationsTable,
	})
	if err != nil {
		return err
	}

	return Migrate(db, dbDriver)
}
