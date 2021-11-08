package telemetry

import (
	"database/sql"
	"log"
	"time"
)

type Aggregator struct {
	DB *sql.DB
}

func NewAggregator(db *sql.DB) *Aggregator {
	return &Aggregator{
		DB: db,
	}
}

func (a *Aggregator) Run(d time.Duration) {
	log.Printf("started aggregator for %s\n", d)
	// Define the duration starts and end.
	// Allow a buffer of the duration to define the start and end.
	// This is to ensure we wait for people not being connected or if they received messages with delay
	runAt := time.Now()
	endsAt := runAt.Add(-d)
	startsAt := endsAt.Add(-d)

	// Query all received message for a specific duration
	receivedMessages, err := queryReceivedMessagesBetween(a.DB, startsAt, endsAt)
	if err != nil {
		log.Fatalf("could not query received message: %s", err)
	}

	// Collect all key uids
	receiverKeyUIDs := make(map[string]bool)
	for _, receivedMessage := range receivedMessages {
		receiverKeyUIDs[receivedMessage.ReceiverKeyUID] = true
	}

	// Ensure the specific key uids received a message after the end of the duration
	// That way we know that this specific key uid has been connected
	for receiverKeyUID := range receiverKeyUIDs {
		ok, err := didReceivedMessageAfter(a.DB, receiverKeyUID, endsAt)
		if err != nil {
			log.Fatalf("could not check key UID: %s, because of %s", receiverKeyUID, err)
		}
		if !ok {
			receiverKeyUIDs[receiverKeyUID] = false
		}
	}

	// Group the received messages by chat id and key uid
	groupedMessages := make(map[string]map[string]int)
	for _, receivedMessage := range receivedMessages {
		// Skip receiver key uid if it has not been connected
		if !receiverKeyUIDs[receivedMessage.ReceiverKeyUID] {
			continue
		}

		if _, ok := groupedMessages[receivedMessage.ChatID]; !ok {
			groupedMessages[receivedMessage.ChatID] = make(map[string]int)
		}
		groupedMessages[receivedMessage.ChatID][receivedMessage.ReceiverKeyUID] += 1
	}

	if len(groupedMessages) == 0 {
		log.Println("no record found, finishing early")
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
			log.Fatalf("could not store received message aggregated: %s", err)
		}
	}
	log.Printf("stored %d chat id records", len(rChatID))

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
		log.Fatalf("could not store received message aggregated: %s", err)
	}
	log.Printf("finished aggregator for %s\n", d)
}
