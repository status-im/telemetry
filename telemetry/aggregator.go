package telemetry

import (
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type Aggregator struct {
	DB     *sql.DB
	logger zap.Logger
}

func NewAggregator(db *sql.DB, logger *zap.Logger) (*Aggregator, error) {
	return &Aggregator{
		DB:     db,
		logger: *logger,
	}, nil
}

func (a *Aggregator) Run(d time.Duration) {
	a.logger.Info("started aggregator", zap.Duration("duration", d))
	// Define the duration starts and end.
	// Allow a buffer of the duration to define the start and end.
	// This is to ensure we wait for people not being connected or if they received messages with delay
	runAt := time.Now()
	endsAt := runAt.Add(-d)
	startsAt := endsAt.Add(-d)

	// Query all received message for a specific duration
	receivedMessages, err := queryReceivedMessagesBetween(a.DB, startsAt, endsAt)
	if err != nil {
		a.logger.Fatal("could not query received message", zap.Error(err))
	}

	// Group the received messages by chat id and key uid
	groupedMessages := make(map[string]map[string]int)
	for _, receivedMessage := range receivedMessages {
		// Skip receiver key uid if it has not been connected or was not in the chat after and before
		ok, err := didReceivedMessageBeforeAndAfterInChat(
			a.DB, receivedMessage.ReceiverKeyUID,
			startsAt,
			endsAt,
			receivedMessage.ChatID,
		)
		if err != nil {
			a.logger.Fatal("could not check message", zap.Int("messageID", receivedMessage.ID), zap.Error(err))
		}
		if !ok {
			continue
		}

		if _, ok := groupedMessages[receivedMessage.ChatID]; !ok {
			groupedMessages[receivedMessage.ChatID] = make(map[string]int)
		}
		groupedMessages[receivedMessage.ChatID][receivedMessage.ReceiverKeyUID] += 1
	}

	if len(groupedMessages) == 0 {
		a.logger.Info("no record found, finishing early")
		return
	}

	// Rch = 1 - count_of_message_missing / total_number_of_messages

	// Calculate the reliability for each channel as:
	// Rch = 1 - count_of_message_missing / total_number_of_messages
	rChatID := make(map[string]float64)
	for chatID, countByKeyUID := range groupedMessages {
		messageMissing := 0
		totalMessages := 0

		max := 0
		for _, count := range countByKeyUID {
			if count > max {
				max = count
			}
		}

		for _, count := range countByKeyUID {
			totalMessages += count
			messageMissing += max - count
		}

		rChatID[chatID] = 1 - float64(messageMissing)/float64(totalMessages)
	}

	// Store all aggregation
	for ChatID, rChatID := range rChatID {
		rma := ReceivedMessageAggregated{
			ChatID:            ChatID,
			DurationInSeconds: int64(d.Seconds()),
			Value:             rChatID,
			RunAt:             runAt.Unix(),
		}
		err := rma.put(a.DB)
		if err != nil {
			a.logger.Fatal("could not store received message aggregated", zap.Error(err))
		}
	}
	a.logger.Sugar().Infof("stored %d chat id records", len(rChatID))

	// Calculate the global reliability R = (R(0) + R(1)+ .... + R(n)) / len(Rch)
	rChatIDTotal := 0.0
	for _, v := range rChatID {
		rChatIDTotal += v
	}

	r := rChatIDTotal / float64(len(rChatID))
	rma := ReceivedMessageAggregated{
		ChatID:            "",
		DurationInSeconds: int64(d.Seconds()),
		Value:             r,
		RunAt:             runAt.Unix(),
	}
	err = rma.put(a.DB)
	if err != nil {
		a.logger.Fatal("could not store received message aggregateds", zap.Error(err))
	}
	a.logger.Info("finished aggregator", zap.Duration("duration", d))
}
