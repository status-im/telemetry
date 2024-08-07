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
	"github.com/status-im/dev-telemetry/pkg/types"
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
	server.Router.HandleFunc("/update-envelope", server.updateEnvelope).Methods("POST")
	server.Router.HandleFunc("/record-metrics", server.createTelemetryData).Methods("POST")

	server.Router.HandleFunc("/waku-metrics", server.createWakuTelemetry).Methods("POST")
	server.Router.HandleFunc("/health", handleHealthCheck).Methods("GET")
	server.Router.Use(server.rateLimit)

	return server
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func (s *Server) createTelemetryData(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var telemetryData []types.TelemetryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&telemetryData); err != nil {
		log.Println(err)
		http.Error(w, "Failed to decode telemetry data", http.StatusBadRequest)
		return
	}

	var errorDetails MetricErrors

	for _, data := range telemetryData {
		switch data.TelemetryType {
		case types.ProtocolStatsMetric:
			var stats ProtocolStats
			err := stats.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.ReceivedEnvelopeMetric:
			var envelope ReceivedEnvelope
			err := envelope.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.SentEnvelopeMetric:
			var envelope SentEnvelope
			err := envelope.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.ErrorSendingEnvelopeMetric:
			var envelopeError ErrorSendingEnvelope
			err := envelopeError.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.PeerCountMetric:
			var peerCount PeerCount
			err := peerCount.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.PeerConnFailureMetric:
			var peerConnFailure PeerConnFailure
			err := peerConnFailure.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		case types.ReceivedMessagesMetric:
			var receivedMessages ReceivedMessage
			err := receivedMessages.process(s.DB, &errorDetails, &data)
			if err != nil {
				continue
			}
		default:
			errorDetails.Append(data.Id, fmt.Sprintf("Unknown telemetry type: %s", data.TelemetryType))
		}
	}

	if errorDetails.Len() > 0 {
		log.Printf("Errors encountered: %v", errorDetails.Get())
	}

	err := respondWithJSON(w, http.StatusCreated, errorDetails.Get())
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

func (s *Server) updateEnvelope(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedEnvelope ReceivedEnvelope
	decoder := json.NewDecoder(r.Body)
	s.logger.Info("update envelope")
	if err := decoder.Decode(&receivedEnvelope.data); err != nil {
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

func (s *Server) createProtocolStats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var protocolStats ProtocolStats
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&protocolStats.data); err != nil {
		s.logger.Error("failed to decode protocol stats", zap.Error(err))

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	peerIDHash := sha256.Sum256([]byte(protocolStats.data.PeerID))
	protocolStats.data.PeerID = hex.EncodeToString(peerIDHash[:])

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

func (s *Server) createWakuTelemetry(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var telemetryData []WakuTelemetryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&telemetryData); err != nil {
		log.Println(err)
		http.Error(w, "Failed to decode telemetry data", http.StatusBadRequest)
		return
	}

	var errorDetails MetricErrors

	for _, data := range telemetryData {
		switch data.TelemetryType {
		case LightPushFilter:
			var pushFilter TelemetryPushFilter
			if err := json.Unmarshal(*data.TelemetryData, &pushFilter); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error decoding lightpush/filter metric: %v", err))
				continue
			}
			if err := pushFilter.put(s.DB); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error saving lightpush/filter metric: %v", err))
				continue
			}
		case LightPushError:
			var pushError TelemetryPushError
			if err := json.Unmarshal(*data.TelemetryData, &pushError); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error decoding lightpush error metric: %v", err))
				continue
			}
			if err := pushError.put(s.DB); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error saving lightpush error metric: %v", err))
				continue
			}
		case Generic:
			var pushGeneric TelemetryGeneric
			if err := json.Unmarshal(*data.TelemetryData, &pushGeneric); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error decoding lightpush generic metric: %v", err))
				continue
			}
			if err := pushGeneric.put(s.DB); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error saving lightpush generic metric: %v", err))
				continue
			}
		default:
			errorDetails.Append(data.Id, fmt.Sprintf("Unknown waku telemetry type: %s", data.TelemetryType))
		}
	}

	if errorDetails.Len() > 0 {
		log.Printf("Errors encountered: %v", errorDetails.Get())
		errorDetailsJSON, err := json.Marshal(errorDetails.Get())
		if err != nil {
			s.logger.Error("failed to marshal error details", zap.Error(err))
			http.Error(w, "Failed to process error details", http.StatusInternalServerError)
			return
		}
		err = respondWithError(w, http.StatusInternalServerError, string(errorDetailsJSON))
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err := respondWithJSON(w, http.StatusCreated, errorDetails.Get())
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

func (s *Server) Start(port int) {
	s.logger.Info("Starting server", zap.Int("port", port))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
