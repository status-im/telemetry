package telemetry

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/status-im/telemetry/lib/common"
	"github.com/status-im/telemetry/lib/metrics"
	"github.com/status-im/telemetry/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	rateLimit        = rate.Limit(10)
	burst            = 1
	cleanUpExecution = 1 * time.Hour * 24
)

type Server struct {
	Router           *mux.Router
	DB               *sql.DB
	logger           *zap.Logger
	rateLimiter      RateLimiter
	ctx              context.Context
	metricsRetention time.Duration
	metrics          map[types.TelemetryType]common.MetricProcessor
}

func NewServer(db *sql.DB, logger *zap.Logger, metricsRetention time.Duration) *Server {
	ctx := context.Background()
	server := &Server{
		Router:           mux.NewRouter().StrictSlash(true),
		DB:               db,
		logger:           logger,
		rateLimiter:      *NewRateLimiter(ctx, rateLimit, burst, logger),
		ctx:              ctx,
		metricsRetention: metricsRetention,
		metrics:          make(map[types.TelemetryType]common.MetricProcessor),
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

func (s *Server) RegisterMetric(t types.TelemetryType, metric common.MetricProcessor) {
	s.metrics[t] = metric
}

func (s *Server) cleanup() {
	if s.metricsRetention == 0 {
		s.logger.Info("Retention set to 0, exiting the cleanup loop")
		return
	}

	timer := time.NewTimer(cleanUpExecution)
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			for name, metric := range s.metrics {
				rows, err := metric.Clean(s.DB, int64(time.Now().Add(-s.metricsRetention).Unix()))
				if err != nil {
					s.logger.Error("Failed to clean up a metric", zap.String("metric", string(name)), zap.Error(err))
					continue
				}
				s.logger.Info("Cleaned up a metric", zap.String("metric", string(name)), zap.Int64("removed", rows))
			}
			timer.Reset(cleanUpExecution)
		}
	}
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

	errorDetails := common.NewMetricErrors(s.logger)

	for _, data := range telemetryData {
		metric, ok := s.metrics[data.TelemetryType]
		if !ok {
			s.logger.Info(fmt.Sprintf("Unknown telemetry type: %s", data.TelemetryType))
			continue
		}

		err := metric.Process(s.ctx, s.DB, errorDetails, &data)
		if err != nil {
			continue
		}
	}

	err := common.RespondWithJSON(w, http.StatusCreated, errorDetails.Get())
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
	var receivedEnvelope metrics.ReceivedEnvelope
	decoder := json.NewDecoder(r.Body)
	s.logger.Info("update envelope")
	if err := decoder.Decode(&receivedEnvelope); err != nil {
		s.logger.Error("failed to decode envelope", zap.Error(err))

		err := common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	err := receivedEnvelope.UpdateProcessingError(s.DB)
	if err != nil {
		s.logger.Error("could not update envelope", zap.Error(err), zap.Any("envelope", receivedEnvelope))
		err := common.RespondWithError(w, http.StatusBadRequest, "Could not update the envelope")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err = common.RespondWithJSON(w, http.StatusCreated, receivedEnvelope)
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
	var protocolStats metrics.ProtocolStats
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&protocolStats); err != nil {
		s.logger.Error("failed to decode protocol stats", zap.Error(err))

		err := common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}
	defer r.Body.Close()

	if err := protocolStats.Put(s.DB); err != nil {
		s.logger.Error("failed to save protocol stats", zap.Error(err))
		err := common.RespondWithError(w, http.StatusInternalServerError, "Could not save protocol stats")
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err := common.RespondWithJSON(w, http.StatusCreated, common.ErrorDetail{})
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
	var telemetryData []metrics.WakuTelemetryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&telemetryData); err != nil {
		log.Println(err)
		http.Error(w, "Failed to decode telemetry data", http.StatusBadRequest)
		return
	}

	errorDetails := common.NewMetricErrors(s.logger)

	for _, data := range telemetryData {
		switch data.TelemetryType {
		case metrics.LightPushFilter:
			var pushFilter metrics.TelemetryPushFilter
			if err := json.Unmarshal(*data.TelemetryData, &pushFilter); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error decoding lightpush/filter metric: %v", err))
				continue
			}
			if err := pushFilter.Put(s.DB); err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					errorDetails.Append(data.Id, "Error saving lightpush/filter metric: Duplicate key value violates unique constraint")
					continue
				}
				errorDetails.Append(data.Id, fmt.Sprintf("Error saving lightpush/filter metric: %v", err))
				continue
			}
		case metrics.Generic:
			var pushGeneric metrics.TelemetryGeneric
			if err := json.Unmarshal(*data.TelemetryData, &pushGeneric); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error decoding lightpush generic metric: %v", err))
				continue
			}
			if err := pushGeneric.Put(s.DB); err != nil {
				errorDetails.Append(data.Id, fmt.Sprintf("Error saving lightpush generic metric: %v", err))
				continue
			}
		default:
			errorDetails.Append(data.Id, fmt.Sprintf("Unknown waku telemetry type: %s", data.TelemetryType))
		}
	}

	if errorDetails.Len() > 0 {
		errorDetailsJSON, err := json.Marshal(errorDetails.Get())
		if err != nil {
			s.logger.Error("failed to marshal error details", zap.Error(err))
			http.Error(w, "Failed to process error details", http.StatusInternalServerError)
			return
		}
		err = common.RespondWithError(w, http.StatusInternalServerError, string(errorDetailsJSON))
		if err != nil {
			s.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	err := common.RespondWithJSON(w, http.StatusCreated, errorDetails.Get())
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

	go s.cleanup()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
