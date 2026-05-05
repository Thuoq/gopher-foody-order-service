package domain

import "time"

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null;index"`
	FoodID    string    `json:"food_id" gorm:"not null"`
	
	// Denormalized fields to capture state at time of order
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	
	Quantity  int       `json:"quantity" gorm:"not null"`
	Subtotal  float64   `json:"subtotal" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
