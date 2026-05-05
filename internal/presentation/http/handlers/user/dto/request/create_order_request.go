package request

type CreateOrderRequest struct {
	RestaurantID    string      `json:"restaurant_id" binding:"required"`
	ShippingAddress string      `json:"shipping_address" binding:"required"`
	Note            string      `json:"note"`
	Items           []OrderItem `json:"items" binding:"required,min=1"`
}

type OrderItem struct {
	FoodID   string `json:"food_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}
