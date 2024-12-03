package main

import (
	"flag"
	"log"
	"time"

	"github.com/status-im/telemetry/lib/database"
	"github.com/status-im/telemetry/lib/metrics"
	"github.com/status-im/telemetry/pkg/types"
	"github.com/status-im/telemetry/telemetry"
	"go.uber.org/zap"

	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	port := flag.Int("port", 8080, "Port number")
	dataSourceName := flag.String("data-source-name", "", "DB URL")
	retention := flag.Duration("retention", 0, "Duration of metrics retention")

	flag.Parse()

	db := database.OpenDb(*dataSourceName, logger)
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

	server := telemetry.NewServer(db, logger, *retention)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"https://lab.waku.org", "https://buddybook.fun"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	server.Router.Use(corsHandler.Handler)

	server.RegisterMetric(types.PeerCountMetric, &metrics.PeerCount{})
	server.RegisterMetric(types.ErrorSendingEnvelopeMetric, &metrics.ErrorSendingEnvelope{})
	server.RegisterMetric(types.PeerConnFailureMetric, &metrics.PeerConnFailure{})
	server.RegisterMetric(types.ProtocolStatsMetric, &metrics.ProtocolStats{})
	server.RegisterMetric(types.ReceivedEnvelopeMetric, &metrics.ReceivedEnvelope{})
	server.RegisterMetric(types.ReceivedMessagesMetric, &metrics.ReceivedMessage{})
	server.RegisterMetric(types.SentEnvelopeMetric, &metrics.SentEnvelope{})
	server.RegisterMetric(types.PeerCountByShardMetric, &metrics.PeerCountByShard{})
	server.RegisterMetric(types.PeerCountByOriginMetric, &metrics.PeerCountByOrigin{})
	server.RegisterMetric(types.MissedMessageMetric, &metrics.MissedMessages{})
	server.RegisterMetric(types.MissedRelevantMessageMetric, &metrics.MissedRelevantMessages{})
	server.RegisterMetric(types.MessageDeliveryConfirmedMetric, &metrics.MessageDeliveryConfirmed{})
	server.RegisterMetric(types.MessageCheckSuccessMetric, &metrics.MessageCheckSuccess{})
	server.RegisterMetric(types.MessageCheckFailureMetric, &metrics.MessageCheckFailure{})
	server.RegisterMetric(types.DialFailureMetric, &metrics.DialFailure{})
	server.RegisterMetric(types.SentMessageTotalMetric, &metrics.SentMessageTotalMetric{})
	server.Start(*port)
}
