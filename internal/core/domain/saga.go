package domain

type SagaStatus string

const (
	SagaStatusStarted     SagaStatus = "started"
	SagaStatusProcessing  SagaStatus = "processing"
	SagaStatusCompensating SagaStatus = "compensating"
	SagaStatusCompleted    SagaStatus = "completed"
	SagaStatusFailed       SagaStatus = "failed"
)

// Event types for the Saga
const (
	EventOrderCreated             = "order.created"
	EventRestaurantValidated      = "restaurant.validated"
	EventRestaurantValidationFailed = "restaurant.validation_failed"
	EventPaymentProcessed         = "payment.processed"
	EventPaymentFailed            = "payment.failed"
)

type SagaEvent struct {
	OrderID      string      `json:"order_id"`
	EventType    string      `json:"event_type"`
	Payload      interface{} `json:"payload"`
	Timestamp    int64       `json:"timestamp"`
}
