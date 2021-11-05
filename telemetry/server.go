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

	return server
}

func (s *Server) createReceivedMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var receivedMessage ReceivedMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receivedMessage); err != nil {
		log.Println(err)

		err := respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	if err := receivedMessage.put(s.DB); err != nil {
		log.Println(err)

		err := respondWithError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err := respondWithJSON(w, http.StatusCreated, receivedMessage)
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
