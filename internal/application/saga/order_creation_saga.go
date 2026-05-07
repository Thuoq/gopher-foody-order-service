package saga

import (
	"context"
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"
	"time"
)

type orderCreationSaga struct {
	orderRepo ports.IOrderRepository
	publisher ports.IMessagePublisher
}

func NewOrderCreationSaga(
	orderRepo ports.IOrderRepository,
	publisher ports.IMessagePublisher,
) ports.ISagaOrchestrator {
	return &orderCreationSaga{
		orderRepo: orderRepo,
		publisher: publisher,
	}
}

func (s *orderCreationSaga) CreateStartEvent(order *domain.Order) *domain.SagaEvent {
	return &domain.SagaEvent{
		OrderID:   order.PublicID,
		EventType: domain.EventOrderCreated,
		Payload: map[string]interface{}{
			"restaurant_id": order.RestaurantID,
			"items":         order.Items,
		},
		Timestamp: time.Now().Unix(),
	}
}

func (s *orderCreationSaga) HandleEvent(ctx context.Context, event domain.SagaEvent) error {
	switch event.EventType {
	case domain.EventRestaurantValidated:
		return s.handleRestaurantValidated(ctx, event)
	case domain.EventRestaurantValidationFailed:
		return s.handleRestaurantValidationFailed(ctx, event)
	case domain.EventPaymentProcessed:
		return s.handlePaymentProcessed(ctx, event)
	case domain.EventPaymentFailed:
		return s.handlePaymentFailed(ctx, event)
	}
	return nil
}

func (s *orderCreationSaga) handleRestaurantValidated(ctx context.Context, event domain.SagaEvent) error {
	// Step 2: Trigger Payment
	// In a real scenario, we would fetch the order to get the total price
	order, err := s.orderRepo.GetByPublicID(ctx, event.OrderID)
	if err != nil {
		return err
	}

	nextEvent := domain.SagaEvent{
		OrderID:   order.PublicID,
		EventType: "payment.process", // Internal command to trigger payment
		Payload: map[string]interface{}{
			"user_id":     order.UserID,
			"amount":      order.TotalPrice,
			"description": "Payment for order " + order.PublicID,
		},
		Timestamp: time.Now().Unix(),
	}

	return s.publisher.Publish(ctx, "payment.process", nextEvent)
}

func (s *orderCreationSaga) handleRestaurantValidationFailed(ctx context.Context, event domain.SagaEvent) error {
	// Compensation: Cancel Order
	return s.orderRepo.UpdateStatusByPublicID(ctx, event.OrderID, domain.StatusCancelled, "Restaurant validation failed")
}

func (s *orderCreationSaga) handlePaymentProcessed(ctx context.Context, event domain.SagaEvent) error {
	// Final Step: Confirm Order
	return s.orderRepo.UpdateStatusByPublicID(ctx, event.OrderID, domain.StatusConfirmed, "Payment successful")
}

func (s *orderCreationSaga) handlePaymentFailed(ctx context.Context, event domain.SagaEvent) error {
	// Compensation: Cancel Order
	return s.orderRepo.UpdateStatusByPublicID(ctx, event.OrderID, domain.StatusCancelled, "Payment failed")
}
