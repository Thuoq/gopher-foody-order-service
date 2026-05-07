package order

import (
	"context"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"

	"github.com/google/uuid"
)

type createOrderUseCase struct {
	orderRepo        ports.IOrderRepository
	restaurantClient ports.RestaurantServiceClient
	sagaOrchestrator ports.ISagaOrchestrator
}

func NewCreateOrderUseCase(
	orderRepo ports.IOrderRepository,
	restaurantClient ports.RestaurantServiceClient,
	sagaOrchestrator ports.ISagaOrchestrator,
) ports.ICreateOrderUseCase {
	return &createOrderUseCase{
		orderRepo:        orderRepo,
		restaurantClient: restaurantClient,
		sagaOrchestrator: sagaOrchestrator,
	}
}

func (uc *createOrderUseCase) Execute(ctx context.Context, input ports.CreateOrderInput) (*domain.Order, error) {
	// 1. Get Food Details via gRPC
	foodIDs := make([]string, len(input.Items))
	itemMap := make(map[string]int)
	for i, item := range input.Items {
		foodIDs[i] = item.FoodID
		itemMap[item.FoodID] = item.Quantity
	}

	foodInfos, err := uc.restaurantClient.GetFoodsInfo(ctx, foodIDs)
	if err != nil {
		return nil, err
	}

	// 2. Build Order Items & Calculate Total
	var totalPrice float64
	orderItems := make([]domain.OrderItem, 0)
	for _, info := range foodInfos {
		qty := itemMap[info.PublicID]
		subtotal := info.Price * float64(qty)
		totalPrice += subtotal

		orderItems = append(orderItems, domain.OrderItem{
			FoodID:   info.PublicID,
			Name:     info.Name,
			Price:    info.Price,
			Quantity: qty,
			Subtotal: subtotal,
		})
	}

	if len(orderItems) == 0 {
		return nil, domain.ErrEmptyOrder
	}

	// 3. Create Order Aggregate
	order := &domain.Order{
		PublicID:        uuid.New().String(),
		UserID:          input.UserID,
		RestaurantID:    input.RestaurantID,
		Status:          domain.StatusPending,
		TotalPrice:      totalPrice,
		ShippingAddress: input.ShippingAddress,
		PaymentStatus:   "unpaid",
		Note:            input.Note,
		Items:           orderItems,
		History: []domain.OrderHistory{
			{
				ToStatus:  domain.StatusPending,
				ChangedBy: input.UserID,
				Reason:    "Order created",
			},
		},
	}

	// 4. Create the Saga Start Event
	event := uc.sagaOrchestrator.CreateStartEvent(order)

	// 5. Save to Repository (Transaction handled by Repo - Order + Outbox)
	if err := uc.orderRepo.Create(ctx, order, event); err != nil {
		return nil, err
	}

	return order, nil
}
