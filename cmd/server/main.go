package main

import (
	"flag"

	"github.com/status-im/dev-telemetry/telemetry"
)

func main() {
	port := flag.Int("port", 8080, "Port number")
	dbUsername := flag.String("db-username", "", "Db username")
	dbPassword := flag.String("db-password", "", "Db password")
	dbName := flag.String("db-name", "", "Db name")

	flag.Parse()

	db := telemetry.OpenDb(
		*dbUsername,
		*dbPassword,
		*dbName,
	)
	defer db.Close()

	server := telemetry.NewServer(db)
	server.Start(*port)
}
