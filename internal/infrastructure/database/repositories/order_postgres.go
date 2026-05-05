package repositories

import (
	"context"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"

	"gorm.io/gorm"
)

type orderPostgresRepository struct {
	db *gorm.DB
}

func NewOrderPostgresRepository(db *gorm.DB) ports.IOrderRepository {
	return &orderPostgresRepository{
		db: db,
	}
}

func (r *orderPostgresRepository) Create(ctx context.Context, order *domain.Order) error {
	// GORM handles the transaction automatically if we pass the whole aggregate
	// and have the correct foreign keys/associations set up.
	// However, it's safer to be explicit in microservices.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *orderPostgresRepository) GetByPublicID(ctx context.Context, publicID string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("History").
		Where("public_id = ?", publicID).
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderPostgresRepository) UpdateStatus(ctx context.Context, orderID uint, newStatus domain.OrderStatus, changedBy string, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Get current status for history
		var currentStatus domain.OrderStatus
		if err := tx.Model(&domain.Order{}).Where("id = ?", orderID).Select("status").Scan(&currentStatus).Error; err != nil {
			return err
		}

		// 2. Update status
		if err := tx.Model(&domain.Order{}).Where("id = ?", orderID).Update("status", newStatus).Error; err != nil {
			return err
		}

		// 3. Create history record
		history := &domain.OrderHistory{
			OrderID:    orderID,
			FromStatus: currentStatus,
			ToStatus:   newStatus,
			Reason:     reason,
			ChangedBy:  changedBy,
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderPostgresRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Items").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}
