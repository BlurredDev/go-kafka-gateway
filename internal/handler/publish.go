package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Sender interface {
	Send(msg []byte, correlationID string) error
}

func MakePublishHandler(sender Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if !json.Valid(body) {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		correlationID := r.Header.Get("X-Correlation-Id")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		if err := sender.Send(body, correlationID); err != nil {
			log.Printf("Failed to send to Kafka: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
