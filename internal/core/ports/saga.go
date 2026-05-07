package ports

import (
	"context"
	"gopher-order-service/internal/core/domain"
)

type IMessagePublisher interface {
	Publish(ctx context.Context, topic string, event domain.SagaEvent) error
}

type ISagaOrchestrator interface {
	CreateStartEvent(order *domain.Order) *domain.SagaEvent
	HandleEvent(ctx context.Context, event domain.SagaEvent) error
}
