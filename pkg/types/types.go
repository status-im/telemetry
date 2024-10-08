package types

import "encoding/json"

type TelemetryType string

const (
	ProtocolStatsMetric            TelemetryType = "ProtocolStats"
	ReceivedEnvelopeMetric         TelemetryType = "ReceivedEnvelope"
	SentEnvelopeMetric             TelemetryType = "SentEnvelope"
	UpdateEnvelopeMetric           TelemetryType = "UpdateEnvelope"
	ReceivedMessagesMetric         TelemetryType = "ReceivedMessages"
	ErrorSendingEnvelopeMetric     TelemetryType = "ErrorSendingEnvelope"
	PeerCountMetric                TelemetryType = "PeerCount"
	PeerConnFailureMetric          TelemetryType = "PeerConnFailure"
	PeerCountByShardMetric         TelemetryType = "PeerCountByShard"
	PeerCountByOriginMetric        TelemetryType = "PeerCountByOrigin"
	MessageCheckSuccessMetric      TelemetryType = "MessageCheckSuccess"
	MessageCheckFailureMetric      TelemetryType = "MessageCheckFailure"
	DialFailureMetric              TelemetryType = "DialFailure"
	StoreConfrimationErrorMetric   TelemetryType = "StoreConfrimationError"
	MissedMessageMetric            TelemetryType = "MissedMessages"
	MissedRelevantMessageMetric    TelemetryType = "MissedRelevantMessages"
	MessageDeliveryConfirmedMetric TelemetryType = "MessageDeliveryConfirmed"
)

type Origin int64

const (
	Unknown Origin = iota
	Discv5
	Static
	PeerExchange
	DNSDiscovery
	Rendezvous
	PeerManager
)

type DialErrorType int

const (
	ErrorUnknown DialErrorType = iota
	ErrorIOTimeout
	ErrorConnectionRefused
	ErrorRelayCircuitFailed
	ErrorRelayNoReservation
	ErrorSecurityNegotiationFailed
	ErrorConcurrentDialSucceeded
	ErrorConcurrentDialFailed
)

type TelemetryRequest struct {
	ID            int              `json:"id"`
	TelemetryType TelemetryType    `json:"telemetry_type"`
	TelemetryData *json.RawMessage `json:"telemetry_data"`
}

type TelemetryRecord struct {
	NodeName      string `json:"nodeName"`
	PeerID        string `json:"peerId"`
	StatusVersion string `json:"statusVersion"`
	DeviceType    string `json:"deviceType"`
}

type PeerCount struct {
	TelemetryRecord
	ID         int    `json:"id"`
	NodeKeyUID string `json:"nodeKeyUID"`
	PeerCount  int    `json:"peerCount"`
	Timestamp  int64  `json:"timestamp"`
}

type PeerConnFailure struct {
	TelemetryRecord
	ID           int    `json:"id"`
	NodeKeyUID   string `json:"nodeKeyUID"`
	FailedPeerId string `json:"failedPeerId"`
	FailureCount int    `json:"failureCount"`
	Timestamp    int64  `json:"timestamp"`
}

type SentEnvelope struct {
	TelemetryRecord
	ID              int    `json:"id"`
	MessageHash     string `json:"messageHash"`
	SentAt          int64  `json:"sentAt"`
	PubsubTopic     string `json:"pubsubTopic"`
	Topic           string `json:"topic"`
	SenderKeyUID    string `json:"senderKeyUID"`
	ProcessingError string `json:"processingError"`
	PublishMethod   string `json:"publishMethod"`
}

type ErrorSendingEnvelope struct {
	ID           int          `json:"id"`
	Error        string       `json:"error"`
	SentEnvelope SentEnvelope `json:"sentEnvelope"`
}

type ReceivedEnvelope struct {
	TelemetryRecord
	ID              int    `json:"id"`
	SentAt          int64  `json:"sentAt"`
	MessageHash     string `json:"messageHash"`
	PubsubTopic     string `json:"pubsubTopic"`
	Topic           string `json:"topic"`
	ReceiverKeyUID  string `json:"receiverKeyUID"`
	ProcessingError string `json:"processingError"`
}

type Metric struct {
	TotalIn  int64   `json:"totalIn"`
	TotalOut int64   `json:"totalOut"`
	RateIn   float64 `json:"rateIn"`
	RateOut  float64 `json:"rateOut"`
}

type ProtocolStats struct {
	PeerID          string `json:"hostID"`
	Relay           Metric `json:"relay"`
	Store           Metric `json:"store"`
	FilterPush      Metric `json:"filter-push"`
	FilterSubscribe Metric `json:"filter-subscribe"`
	Lightpush       Metric `json:"lightpush"`
}

type ReceivedMessage struct {
	TelemetryRecord
	ID             int    `json:"id"`
	ChatID         string `json:"chatId"`
	MessageHash    string `json:"messageHash"`
	MessageID      string `json:"messageId"`
	MessageType    string `json:"messageType"`
	MessageSize    int    `json:"messageSize"`
	ReceiverKeyUID string `json:"receiverKeyUID"`
	SentAt         int64  `json:"sentAt"`
	Topic          string `json:"topic"`
	PubsubTopic    string `json:"pubsubTopic"`
}

type PeerCountByShard struct {
	TelemetryRecord
	ID        int   `json:"id"`
	Count     int   `json:"count"`
	Shard     int   `json:"shard"`
	Timestamp int64 `json:"timestamp"`
}

type PeerCountByOrigin struct {
	TelemetryRecord
	ID        int    `json:"id"`
	Count     int    `json:"count"`
	Origin    Origin `json:"origin"`
	Timestamp int64  `json:"timestamp"`
}

type MessageCheckSuccess struct {
	TelemetryRecord
	MessageHash string `json:"messageHash"`
	Timestamp   int64  `json:"timestamp"`
}

type MessageCheckFailure struct {
	TelemetryRecord
	MessageHash string `json:"messageHash"`
	Timestamp   int64  `json:"timestamp"`
}
type DialFailure struct {
	TelemetryRecord
	ErrorType DialErrorType `json:"errorType"`
	ErrorMsg  string        `json:"errorMsg"`
	Protocols string        `json:"protocols"`
	Timestamp int64         `json:"timestamp"`
}

type MissedMessage struct {
	TelemetryRecord
	ContentTopic string `json:"contentTopic"`
	MessageHash  string `json:"messageHash"`
	SentAt       int64  `json:"sentAt"`
	PubsubTopic  string `json:"pubsubTopic"`
}

type MissedRelevantMessage struct {
	TelemetryRecord
	ContentTopic string `json:"contentTopic"`
	MessageHash  string `json:"messageHash"`
	SentAt       int64  `json:"sentAt"`
	PubsubTopic  string `json:"pubsubTopic"`
}

type MessageDeliveryConfirmed struct {
	TelemetryRecord
	MessageHash string `json:"messageHash"`
	Timestamp   int64  `json:"timestamp"`
}
