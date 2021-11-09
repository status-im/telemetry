package telemetry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	server := &Server{
		Router: mux.NewRouter().StrictSlash(true),
		DB:     db,
	}

	server.Router.HandleFunc("/received-messages", server.createReceivedMessages).Methods("POST")
	server.Router.HandleFunc("/health", handleHealthCheck).Methods("GET")

	return server
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func (s *Server) createReceivedMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedMessages []ReceivedMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receivedMessages); err != nil {
		log.Println(err)

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	var ids []int
	for _, receivedMessage := range receivedMessages {
		if err := receivedMessage.put(s.DB); err != nil {
			log.Println("could not save message", err, receivedMessage)
			continue
		}
		ids = append(ids, receivedMessage.ID)
	}

	if len(ids) != len(receivedMessages) {
		err := respondWithError(w, http.StatusInternalServerError, "Could not save all record")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err := respondWithJSON(w, http.StatusCreated, receivedMessages)
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
	log.Printf("Starting server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
