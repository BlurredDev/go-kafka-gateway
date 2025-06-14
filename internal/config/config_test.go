package config

import (
	"os"
	"testing"
)

func TestLoad_WithEnvVars(t *testing.T) {
	os.Setenv("KAFKA_BROKER", "test-broker:1234")
	os.Setenv("KAFKA_TOPIC", "test-topic")
	os.Setenv("HTTP_ADDR", ":9090")
	defer func() {
		os.Unsetenv("KAFKA_BROKER")
		os.Unsetenv("KAFKA_TOPIC")
		os.Unsetenv("HTTP_ADDR")
	}()

	cfg := Load()
	if cfg.KafkaBroker != "test-broker:1234" {
		t.Errorf("expected KafkaBroker 'test-broker:1234', got %s", cfg.KafkaBroker)
	}
	if cfg.KafkaTopic != "test-topic" {
		t.Errorf("expected KafkaTopic 'test-topic', got %s", cfg.KafkaTopic)
	}
	if cfg.HTTPAddr != ":9090" {
		t.Errorf("expected HTTPAddr ':9090', got %s", cfg.HTTPAddr)
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("KAFKA_BROKER")
	os.Unsetenv("KAFKA_TOPIC")
	os.Unsetenv("HTTP_ADDR")

	cfg := Load()
	if cfg.KafkaBroker != "localhost:9092" {
		t.Errorf("expected default KafkaBroker 'localhost:9092', got %s", cfg.KafkaBroker)
	}
	if cfg.KafkaTopic != "default-topic" {
		t.Errorf("expected default KafkaTopic 'default-topic', got %s", cfg.KafkaTopic)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("expected default HTTPAddr ':8080', got %s", cfg.HTTPAddr)
	}
}
