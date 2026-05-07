package ports

import (
	"context"
	"gopher-order-service/internal/core/domain"
)

type IOutboxRepository interface {
	Create(ctx context.Context, message *domain.OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]domain.OutboxMessage, error)
	MarkProcessed(ctx context.Context, id uint) error
	MarkFailed(ctx context.Context, id uint) error
}
