package messaging

import (
	"context"
	"encoding/json"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"
)

type outboxPublisher struct {
	outboxRepo ports.IOutboxRepository
}

func NewOutboxPublisher(outboxRepo ports.IOutboxRepository) ports.IMessagePublisher {
	return &outboxPublisher{
		outboxRepo: outboxRepo,
	}
}

func (p *outboxPublisher) Publish(ctx context.Context, topic string, event domain.SagaEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outbox := &domain.OutboxMessage{
		Topic:   topic,
		Key:     event.OrderID,
		Payload: payload,
		Status:  domain.OutboxStatusPending,
	}

	return p.outboxRepo.Create(ctx, outbox)
}
