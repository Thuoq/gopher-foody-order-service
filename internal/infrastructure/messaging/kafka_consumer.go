package messaging

import (
	"context"
	"encoding/json"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"
	"strings"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	reader       *kafka.Reader
	orchestrator ports.ISagaOrchestrator
	logger       *zap.Logger
}

func NewKafkaConsumer(brokers string, orchestrator ports.ISagaOrchestrator, logger *zap.Logger) *KafkaConsumer {
	// Topics this orchestrator listens to (responses from other services)
	topics := []string{
		domain.EventRestaurantValidated,
		domain.EventRestaurantValidationFailed,
		domain.EventPaymentProcessed,
		domain.EventPaymentFailed,
	}

	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     strings.Split(brokers, ","),
			GroupID:     "order-saga-worker",
			GroupTopics: topics,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
		}),
		orchestrator: orchestrator,
		logger:       logger,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	c.logger.Info("Saga Worker started, listening to topics...")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("Failed to read message from Kafka", zap.Error(err))
			continue
		}

		var event domain.SagaEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			c.logger.Error("Failed to unmarshal Saga event", zap.Error(err))
			continue
		}

		c.logger.Info("Received Saga event",
			zap.String("topic", m.Topic),
			zap.String("order_id", event.OrderID),
		)

		// Pass event back to the orchestrator
		if err := c.orchestrator.HandleEvent(ctx, event); err != nil {
			c.logger.Error("Orchestrator failed to handle event",
				zap.String("order_id", event.OrderID),
				zap.Error(err),
			)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
