package messaging

import (
	"context"
	"encoding/json"
	"gopher-order-service/internal/core/domain"
	"strings"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaPublisher(brokers string, logger *zap.Logger) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(strings.Split(brokers, ",")...),
			Balancer: &kafka.LeastBytes{},
			Async:    false, // We want to ensure delivery for Saga
		},
		logger: logger,
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, event domain.SagaEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(event.OrderID),
		Value: payload,
	})

	if err != nil {
		p.logger.Error("Failed to publish message to Kafka",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return err
	}

	p.logger.Info("Message published to Kafka",
		zap.String("topic", topic),
		zap.String("order_id", event.OrderID),
	)
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
