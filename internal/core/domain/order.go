package domain

import (
	"time"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusPreparing  OrderStatus = "preparing"
	StatusShipping   OrderStatus = "shipping"
	StatusCompleted  OrderStatus = "completed"
	StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	PublicID        string         `json:"public_id" gorm:"unique;not null;index"`
	UserID          string         `json:"user_id" gorm:"not null;index"`
	RestaurantID    string         `json:"restaurant_id" gorm:"not null;index"`
	Status          OrderStatus    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TotalPrice      float64        `json:"total_price" gorm:"not null"`
	ShippingAddress string         `json:"shipping_address" gorm:"type:text;not null"`
	PaymentStatus   string         `json:"payment_status" gorm:"default:'unpaid'"`
	PaymentMethod   string         `json:"payment_method"`
	Note            string         `json:"note"`
	Items           []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
	History         []OrderHistory `json:"history" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
