package main

import (
	"flag"
	"log"
	"time"

	"github.com/status-im/dev-telemetry/telemetry"
	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	port := flag.Int("port", 8080, "Port number")
	dataSourceName := flag.String("data-source-name", "", "DB URL")

	flag.Parse()

	db := telemetry.OpenDb(*dataSourceName, logger)
	defer db.Close()

	aggregator, err := telemetry.NewAggregator(db, logger)
	if err != nil {
		logger.Fatal("Error creating aggregator", zap.Error(err))
	}

	c := cron.New()
	_, err = c.AddFunc("0 * * * *", func() {
		aggregator.Run(time.Hour)
	})

	if err != nil {
		logger.Fatal("Error adding cron job", zap.Error(err))
	}

	c.Start()
	defer c.Stop()

	server := telemetry.NewServer(db, logger)
	server.Start(*port)
}
