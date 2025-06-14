package config

import "os"

type Config struct {
	KafkaBroker string
	KafkaTopic  string
	HTTPAddr    string
	DLQTopic    string // e.g. from env: KAFKA_DLQ_TOPIC
}

func Load() Config {
	return Config{
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "default-topic"),
		DLQTopic:    getEnv("KAFKA_DLQ_TOPIC", "default-dql-topic"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
