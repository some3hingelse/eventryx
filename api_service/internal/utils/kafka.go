package utils

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

var KafkaProducer Producer

func CreateKafkaProducer(brokers []string, topic string) {
	KafkaProducer = Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *Producer) SendMessageToKafka(ctx context.Context, message []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
