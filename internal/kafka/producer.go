package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type writerInterface interface {
	WriteMessages(context.Context, ...kafka.Message) error
	Close() error
}

type Producer struct {
	mainWriter writerInterface
	dlqWriter  writerInterface
}

func NewProducer(broker, topic, dlqTopic string) *Producer {
	return &Producer{
		mainWriter: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
		dlqWriter: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        dlqTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) Send(message []byte, correlationID string) error {
	msg := kafka.Message{
		Key:   []byte(correlationID),
		Value: message,
		Headers: []kafka.Header{
			{
				Key:   "X-Correlation-Id",
				Value: []byte(correlationID),
			},
		},
	}

	err := p.mainWriter.WriteMessages(context.Background(), msg)
	if err != nil && p.dlqWriter != nil {
		log.Printf("Main Kafka send failed: %v. Writing to DLQ...", err)
		dlqErr := p.dlqWriter.WriteMessages(context.Background(), msg)
		if dlqErr != nil {
			log.Printf("Failed to write to DLQ as well: %v", dlqErr)
		}
		return err
	}

	return err
}

func (p *Producer) Close() error {
	if p.dlqWriter != nil {
		_ = p.dlqWriter.Close()
	}
	return p.mainWriter.Close()
}
