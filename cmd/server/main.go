package main

import (
	"flag"

	"github.com/status-im/dev-telemetry/telemetry"
)

func main() {
	port := flag.Int("port", 8080, "Port number")
	dataSourceName := flag.String("data-source-name", "", "DB URL")

	flag.Parse()

	db := telemetry.OpenDb(*dataSourceName)
	defer db.Close()

	server := telemetry.NewServer(db)
	server.Start(*port)
}
