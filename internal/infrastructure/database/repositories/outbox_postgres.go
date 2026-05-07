package repositories

import (
	"context"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"

	"gorm.io/gorm"
)

type outboxPostgresRepository struct {
	db *gorm.DB
}

func NewOutboxPostgresRepository(db *gorm.DB) ports.IOutboxRepository {
	return &outboxPostgresRepository{
		db: db,
	}
}

func (r *outboxPostgresRepository) Create(ctx context.Context, message *domain.OutboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *outboxPostgresRepository) GetPending(ctx context.Context, limit int) ([]domain.OutboxMessage, error) {
	var messages []domain.OutboxMessage
	err := r.db.WithContext(ctx).
		Where("status = ?", domain.OutboxStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *outboxPostgresRepository) MarkProcessed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.OutboxMessage{}).
		Where("id = ?", id).
		Update("status", domain.OutboxStatusProcessed).Error
}

func (r *outboxPostgresRepository) MarkFailed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.OutboxMessage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      domain.OutboxStatusFailed,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error
}
