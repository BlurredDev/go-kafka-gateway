package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
)

type mockWriter struct {
	calledMessages []kafka.Message
	failWrite      bool
	closed         bool
}

func (m *mockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	m.calledMessages = append(m.calledMessages, msgs...)
	if m.failWrite {
		return errors.New("write failed")
	}
	return nil
}

func (m *mockWriter) Close() error {
	m.closed = true
	return nil
}

func TestProducer_Send_Success(t *testing.T) {
	mainWriter := &mockWriter{}
	dlqWriter := &mockWriter{}
	producer := &Producer{
		mainWriter: mainWriter,
		dlqWriter:  dlqWriter,
	}

	msg := []byte(`{"foo":"bar"}`)
	cid := "cid-123"

	err := producer.Send(msg, cid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mainWriter.calledMessages) != 1 {
		t.Fatalf("expected 1 message sent to mainWriter, got %d", len(mainWriter.calledMessages))
	}

	sentMsg := mainWriter.calledMessages[0]
	if string(sentMsg.Value) != string(msg) {
		t.Errorf("expected message value %s but got %s", string(msg), string(sentMsg.Value))
	}

	if string(sentMsg.Key) != cid {
		t.Errorf("expected key %s but got %s", cid, string(sentMsg.Key))
	}

	if len(sentMsg.Headers) != 1 || sentMsg.Headers[0].Key != "X-Correlation-Id" || string(sentMsg.Headers[0].Value) != cid {
		t.Errorf("expected correlation ID header with key %s and value %s", "X-Correlation-Id", cid)
	}
}

func TestProducer_Send_FallbackToDLQ(t *testing.T) {
	mainWriter := &mockWriter{failWrite: true}
	dlqWriter := &mockWriter{}
	producer := &Producer{
		mainWriter: mainWriter,
		dlqWriter:  dlqWriter,
	}

	msg := []byte(`{"foo":"bar"}`)
	cid := "cid-456"

	err := producer.Send(msg, cid)
	if err == nil {
		t.Fatal("expected error from mainWriter")
	}

	if len(mainWriter.calledMessages) != 1 {
		t.Errorf("expected 1 message sent to mainWriter")
	}

	if len(dlqWriter.calledMessages) != 1 {
		t.Errorf("expected 1 message sent to dlqWriter on fallback")
	}

	sentMsg := dlqWriter.calledMessages[0]
	if string(sentMsg.Value) != string(msg) {
		t.Errorf("expected message value %s but got %s", string(msg), string(sentMsg.Value))
	}

	if string(sentMsg.Key) != cid {
		t.Errorf("expected key %s but got %s", cid, string(sentMsg.Key))
	}
}

func TestProducer_Close(t *testing.T) {
	mainWriter := &mockWriter{}
	dlqWriter := &mockWriter{}
	producer := &Producer{
		mainWriter: mainWriter,
		dlqWriter:  dlqWriter,
	}

	err := producer.Close()
	if err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}

	if !mainWriter.closed {
		t.Error("expected mainWriter to be closed")
	}

	if !dlqWriter.closed {
		t.Error("expected dlqWriter to be closed")
	}
}

func TestNewProducer(t *testing.T) {
	broker := "broker:9092"
	topic := "topic1"
	dlq := "dlq-topic"

	p := NewProducer(broker, topic, dlq)
	if p == nil {
		t.Fatal("expected non-nil Producer")
	}

	// Since mainWriter and dlqWriter are interfaces, check topic by type assertion
	mainWriter, ok := p.mainWriter.(*kafka.Writer)
	if !ok {
		t.Fatal("mainWriter is not *kafka.Writer")
	}
	if mainWriter.Topic != topic {
		t.Errorf("expected mainWriter.Topic %s, got %s", topic, mainWriter.Topic)
	}

	dlqWriter, ok := p.dlqWriter.(*kafka.Writer)
	if !ok {
		t.Fatal("dlqWriter is not *kafka.Writer")
	}
	if dlqWriter.Topic != dlq {
		t.Errorf("expected dlqWriter.Topic %s, got %s", dlq, dlqWriter.Topic)
	}
}

func TestProducer_Send_DLQFailure(t *testing.T) {
	mainWriter := &mockWriter{failWrite: true}
	dlqWriter := &mockWriter{failWrite: true}

	producer := &Producer{
		mainWriter: mainWriter,
		dlqWriter:  dlqWriter,
	}

	err := producer.Send([]byte(`{"key":"value"}`), "correlation-id")

	if err == nil || err.Error() != "write failed" {
		t.Errorf("expected DLQ error, got: %v", err)
	}
}
