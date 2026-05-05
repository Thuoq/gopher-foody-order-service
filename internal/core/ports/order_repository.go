package ports

import (
	"context"
	"gopher-order-service/internal/core/domain"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByPublicID(ctx context.Context, publicID string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, orderID uint, newStatus domain.OrderStatus, changedBy string, reason string) error
	ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Order, int64, error)
}

type CreateOrderInput struct {
	UserID          string
	RestaurantID    string
	ShippingAddress string
	Note            string
	Items           []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	FoodID   string
	Quantity int
}

type ICreateOrderUseCase interface {
	Execute(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
}

type IGetOrderDetailUseCase interface {
	Execute(ctx context.Context, userID string, publicID string) (*domain.Order, error)
}

type FoodInfo struct {
	PublicID string
	Name     string
	Price    float64
}

type RestaurantServiceClient interface {
	GetFoodsInfo(ctx context.Context, foodIDs []string) ([]FoodInfo, error)
}
