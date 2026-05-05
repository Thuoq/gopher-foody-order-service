package domain

import "time"

type OrderHistory struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	OrderID    uint        `json:"order_id" gorm:"not null;index"`
	FromStatus OrderStatus `json:"from_status"`
	ToStatus   OrderStatus `json:"to_status" gorm:"not null"`
	Reason     string      `json:"reason"`
	ChangedBy  string      `json:"changed_by" gorm:"not null"` // UserID or 'system'
	CreatedAt  time.Time   `json:"created_at"`
}
