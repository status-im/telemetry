package main

import (
	"flag"
	"time"

	"github.com/status-im/dev-telemetry/telemetry"
)

func main() {
	seconds := flag.Int("seconds", 3600, "Number of seconds to aggregate")
	dataSourceName := flag.String("data-source-name", "", "DB URL")

	flag.Parse()

	db := telemetry.OpenDb(*dataSourceName)
	defer db.Close()

	aggregator := telemetry.NewAggregator(db)
	aggregator.Run(time.Duration(*seconds))
}
