package main

import (
	"flag"
	"log"
	"time"

	"github.com/status-im/dev-telemetry/telemetry"

	"github.com/robfig/cron/v3"
)

func main() {
	port := flag.Int("port", 8080, "Port number")
	dataSourceName := flag.String("data-source-name", "", "DB URL")

	flag.Parse()

	db := telemetry.OpenDb(*dataSourceName)
	defer db.Close()

	aggregator := telemetry.NewAggregator(db)
	c := cron.New()
	_, err := c.AddFunc("0 * * * *", func() {
		aggregator.Run(time.Hour)
	})

	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	c.Start()
	defer c.Stop()

	server := telemetry.NewServer(db)
	server.Start(*port)
}
