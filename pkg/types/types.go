package types

import "encoding/json"

type TelemetryType string

const (
	ProtocolStatsMetric        TelemetryType = "ProtocolStats"
	ReceivedEnvelopeMetric     TelemetryType = "ReceivedEnvelope"
	SentEnvelopeMetric         TelemetryType = "SentEnvelope"
	UpdateEnvelopeMetric       TelemetryType = "UpdateEnvelope"
	ReceivedMessagesMetric     TelemetryType = "ReceivedMessages"
	ErrorSendingEnvelopeMetric TelemetryType = "ErrorSendingEnvelope"
	PeerCountMetric            TelemetryType = "PeerCount"
	PeerConnFailureMetric      TelemetryType = "PeerConnFailure"
)

type TelemetryRequest struct {
	Id            int              `json:"id"`
	TelemetryType TelemetryType    `json:"telemetry_type"`
	TelemetryData *json.RawMessage `json:"telemetry_data"`
}

type CommonFields struct {
	NodeName      string `json:"nodeName"`
	PeerID        string `json:"peerId"`
	StatusVersion string `json:"statusVersion"`
	DeviceType    string `json:"deviceType"`
}

func (c CommonFields) GetNodeName() string {
	return c.NodeName
}

func (c CommonFields) GetPeerID() string {
	return c.PeerID
}

func (c CommonFields) GetStatusVersion() string {
	return c.StatusVersion
}

func (c CommonFields) GetDeviceType() string {
	return c.DeviceType
}

type CommonFieldsAccessor interface {
	GetNodeName() string
	GetPeerID() string
	GetStatusVersion() string
	GetDeviceType() string
}

type PeerCount struct {
	CommonFields
	ID         int    `json:"id"`
	NodeKeyUid string `json:"nodeKeyUid"`
	PeerCount  int    `json:"peerCount"`
	Timestamp  int64  `json:"timestamp"`
}

type PeerConnFailure struct {
	CommonFields
	ID           int    `json:"id"`
	NodeKeyUid   string `json:"nodeKeyUid"`
	FailedPeerId string `json:"failedPeerId"`
	FailureCount int    `json:"failureCount"`
	Timestamp    int64  `json:"timestamp"`
}

type SentEnvelope struct {
	CommonFields
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
	Id           int          `json:"id"`
	Error        string       `json:"error"`
	SentEnvelope SentEnvelope `json:"sentEnvelope"`
}

type ReceivedEnvelope struct {
	CommonFields
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
	CommonFields
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
