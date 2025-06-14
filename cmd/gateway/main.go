package main

import (
	"log"
	"net/http"

	"github.com/BlurredDev/go-kafka-gateway/internal/config"
	"github.com/BlurredDev/go-kafka-gateway/internal/handler"
	"github.com/BlurredDev/go-kafka-gateway/internal/kafka"
	"github.com/gorilla/mux"
)

func Run(cfg config.Config) error {
	var sender handler.Sender = kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.DLQTopic)
	defer sender.(*kafka.Producer).Close()

	r := mux.NewRouter()
	r.HandleFunc("/publish", handler.MakePublishHandler(sender)).Methods("POST")

	// Add 405 Method Not Allowed handler
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
	})

	log.Printf("Starting gateway on %s", cfg.HTTPAddr)
	return http.ListenAndServe(cfg.HTTPAddr, r)
}

func StartGateway() error {
	cfg := config.Load()
	return Run(cfg)
}

func main() {
	if err := StartGateway(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
