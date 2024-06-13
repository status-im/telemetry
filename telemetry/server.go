package telemetry

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	RATE_LIMIT = rate.Limit(10)
	BURST      = 1
)

type Server struct {
	Router      *mux.Router
	DB          *sql.DB
	logger      *zap.Logger
	rateLimiter RateLimiter
	ctx         context.Context
}

func NewServer(db *sql.DB, logger *zap.Logger) *Server {
	ctx := context.Background()
	server := &Server{
		Router:      mux.NewRouter().StrictSlash(true),
		DB:          db,
		logger:      logger,
		rateLimiter: *NewRateLimiter(ctx, RATE_LIMIT, BURST, logger),
		ctx:         ctx,
	}

	server.Router.HandleFunc("/protocol-stats", server.createProtocolStats).Methods("POST")
	server.Router.HandleFunc("/received-messages", server.createReceivedMessages).Methods("POST")
	server.Router.HandleFunc("/received-envelope", server.createReceivedEnvelope).Methods("POST")
	server.Router.HandleFunc("/sent-envelope", server.createSentEnvelope).Methods("POST")
	server.Router.HandleFunc("/update-envelope", server.updateEnvelope).Methods("POST")
	server.Router.HandleFunc("/health", handleHealthCheck).Methods("GET")
	server.Router.HandleFunc("/record-metrics", server.createTelemetryData).Methods("POST")
	server.Router.Use(server.rateLimit)

	return server
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

type TelemetryType string

const (
	ProtocolStatsMetric    TelemetryType = "ProtocolStats"
	ReceivedEnvelopeMetric TelemetryType = "ReceivedEnvelope"
	SentEnvelopeMetric     TelemetryType = "SentEnvelope"
	UpdateEnvelopeMetric   TelemetryType = "UpdateEnvelope"
	ReceivedMessagesMetric TelemetryType = "ReceivedMessages"
)

type TelemetryRequest struct {
	Id            int              `json:"id"`
	TelemetryType TelemetryType    `json:"telemetry_type"`
	TelemetryData *json.RawMessage `json:"telemetry_data"`
}

func (s *Server) createTelemetryData(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var telemetryData []TelemetryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&telemetryData); err != nil {
		log.Println(err)
		http.Error(w, "Failed to decode telemetry data", http.StatusBadRequest)
		return
	}

	var errorDetails []map[string]interface{}

	for _, data := range telemetryData {
		switch data.TelemetryType {
		case ProtocolStatsMetric:
			var stats ProtocolStats
			if err := json.Unmarshal(*data.TelemetryData, &stats); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error decoding protocol stats: %v", err)})
				continue
			}
			if err := stats.put(s.DB); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error saving protocol stats: %v", err)})
				continue
			}
		case ReceivedEnvelopeMetric:
			var envelope ReceivedEnvelope
			if err := json.Unmarshal(*data.TelemetryData, &envelope); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error decoding received envelope: %v", err)})
				continue
			}
			if err := envelope.put(s.DB); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error saving received envelope: %v", err)})
				continue
			}
		case SentEnvelopeMetric:
			var envelope SentEnvelope
			if err := json.Unmarshal(*data.TelemetryData, &envelope); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error decoding sent envelope: %v", err)})
				continue
			}
			if err := envelope.put(s.DB); err != nil {
				errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Error saving sent envelope: %v", err)})
				continue
			}
		default:
			errorDetails = append(errorDetails, map[string]interface{}{"Id": data.Id, "Error": fmt.Sprintf("Unknown telemetry type: %s", data.TelemetryType)})
		}
	}

	if len(errorDetails) > 0 {
		log.Printf("Errors encountered: %v", errorDetails)
	}

	err := respondWithJSON(w, http.StatusCreated, errorDetails)
	if err != nil {
		log.Println(err)
	}

	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}

func (s *Server) createReceivedMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedMessages []ReceivedMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receivedMessages); err != nil {
		s.logger.Error("failed to decode messages", zap.Error(err))

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	var ids []int
	for _, receivedMessage := range receivedMessages {
		if err := receivedMessage.put(s.DB); err != nil {
			s.logger.Error("could not save message", zap.Error(err), zap.Any("receivedMessage", receivedMessage))
			continue
		}
		ids = append(ids, receivedMessage.ID)
	}

	if len(ids) != len(receivedMessages) {
		err := respondWithError(w, http.StatusInternalServerError, "Could not save all record")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err := respondWithJSON(w, http.StatusCreated, receivedMessages)
	if err != nil {
		s.logger.Error("failed to respond", zap.Error(err))
		return
	}

	s.logger.Info(
		"handled received message",
		zap.String("method", r.Method),
		zap.String("requestURI", r.RequestURI),
		zap.Duration("duration", time.Since(start)),
	)
}

func (s *Server) createReceivedEnvelope(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedEnvelope ReceivedEnvelope
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receivedEnvelope); err != nil {
		s.logger.Error("failed to decode envelope", zap.Error(err))

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
			return
		}
		return
	}
	defer r.Body.Close()

	err := receivedEnvelope.put(s.DB)
	if err != nil {
		s.logger.Error("could not save envelope", zap.Error(err), zap.Any("envelope", receivedEnvelope))
		err := respondWithError(w, http.StatusBadRequest, "Could not save the envelope")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
	}

	err = respondWithJSON(w, http.StatusCreated, receivedEnvelope)
	if err != nil {
		s.logger.Error("failed to respond", zap.Error(err))
		return
	}

	s.logger.Info(
		"handled received envelope",
		zap.String("method", r.Method),
		zap.String("requestURI", r.RequestURI),
		zap.Duration("duration", time.Since(start)),
	)
}

func (s *Server) updateEnvelope(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedEnvelope ReceivedEnvelope
	decoder := json.NewDecoder(r.Body)
	s.logger.Info("update envelope")
	if err := decoder.Decode(&receivedEnvelope); err != nil {
		s.logger.Error("failed to decode envelope", zap.Error(err))

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	err := receivedEnvelope.updateProcessingError(s.DB)
	if err != nil {
		s.logger.Error("could not update envelope", zap.Error(err), zap.Any("envelope", receivedEnvelope))
		err := respondWithError(w, http.StatusBadRequest, "Could not update the envelope")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err = respondWithJSON(w, http.StatusCreated, receivedEnvelope)
	if err != nil {
		s.logger.Error("failed to respond", zap.Error(err))
		return
	}

	s.logger.Info(
		"handled update message",
		zap.String("method", r.Method),
		zap.String("requestURI", r.RequestURI),
		zap.Duration("duration", time.Since(start)),
	)
}

func (s *Server) createSentEnvelope(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var sentEnvelope SentEnvelope
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sentEnvelope); err != nil {
		log.Println(err)

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	err := sentEnvelope.put(s.DB)
	if err != nil {
		log.Println("could not save envelope", err, sentEnvelope)
	}

	err = respondWithJSON(w, http.StatusCreated, sentEnvelope)
	if err != nil {
		log.Println(err)
	}

	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}

func (s *Server) createProtocolStats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var protocolStats ProtocolStats
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&protocolStats); err != nil {
		s.logger.Error("failed to decode protocol stats", zap.Error(err))

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	peerIDHash := sha256.Sum256([]byte(protocolStats.PeerID))
	protocolStats.PeerID = hex.EncodeToString(peerIDHash[:])

	if err := protocolStats.put(s.DB); err != nil {
		s.logger.Error("failed to save protocol stats", zap.Error(err))
		err := respondWithError(w, http.StatusInternalServerError, "Could not save protocol stats")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err := respondWithJSON(w, http.StatusCreated, map[string]string{"error": ""})
	if err != nil {
		s.logger.Error("failed to respond", zap.Error(err))
		return
	}

	s.logger.Info(
		"handled protocol stats",
		zap.String("method", r.Method),
		zap.String("requestURI", r.RequestURI),
		zap.Duration("duration", time.Since(start)),
	)
}

func (s *Server) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		limiter := s.rateLimiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start(port int) {
	s.logger.Info("Starting server", zap.Int("port", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
