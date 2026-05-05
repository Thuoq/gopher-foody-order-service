package response

import "time"

type OrderResponse struct {
	ID              string               `json:"id"`
	UserID          string               `json:"user_id"`
	RestaurantID    string               `json:"restaurant_id"`
	Status          string               `json:"status"`
	TotalPrice      float64              `json:"total_price"`
	ShippingAddress string               `json:"shipping_address"`
	PaymentStatus   string               `json:"payment_status"`
	PaymentMethod   string               `json:"payment_method"`
	Note            string               `json:"note"`
	Items           []OrderItemResponse  `json:"items"`
	History         []HistoryResponse    `json:"history,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
}

type OrderItemResponse struct {
	FoodID   string  `json:"food_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type HistoryResponse struct {
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	Reason     string    `json:"reason"`
	ChangedBy  string    `json:"changed_by"`
	CreatedAt  time.Time `json:"created_at"`
}
