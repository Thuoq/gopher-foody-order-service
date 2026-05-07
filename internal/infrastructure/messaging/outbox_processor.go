package messaging

import (
	"context"
	"encoding/json"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"
	"time"

	"go.uber.org/zap"
)

type OutboxProcessor struct {
	outboxRepo ports.IOutboxRepository
	publisher  ports.IMessagePublisher // This will be the real KafkaPublisher
	logger     *zap.Logger
	interval   time.Duration
}

func NewOutboxProcessor(
	outboxRepo ports.IOutboxRepository,
	publisher ports.IMessagePublisher,
	logger *zap.Logger,
) *OutboxProcessor {
	return &OutboxProcessor{
		outboxRepo: outboxRepo,
		publisher:  publisher,
		logger:     logger,
		interval:   2 * time.Second,
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	p.logger.Info("Outbox Processor started")
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.process(ctx)
		}
	}
}

func (p *OutboxProcessor) process(ctx context.Context) {
	messages, err := p.outboxRepo.GetPending(ctx, 10)
	if err != nil {
		p.logger.Error("Failed to fetch pending outbox messages", zap.Error(err))
		return
	}

	for _, msg := range messages {
		var event domain.SagaEvent
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			p.logger.Error("Failed to unmarshal outbox payload", zap.Uint("id", msg.ID), zap.Error(err))
			p.outboxRepo.MarkFailed(ctx, msg.ID)
			continue
		}

		err = p.publisher.Publish(ctx, msg.Topic, event)
		if err != nil {
			p.logger.Error("Failed to publish outbox message to Kafka", zap.Uint("id", msg.ID), zap.Error(err))
			p.outboxRepo.MarkFailed(ctx, msg.ID)
			continue
		}

		if err := p.outboxRepo.MarkProcessed(ctx, msg.ID); err != nil {
			p.logger.Error("Failed to mark outbox message as processed", zap.Uint("id", msg.ID), zap.Error(err))
		}
	}
}
