package domain

import (
	"encoding/json"
	"time"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusProcessed OutboxStatus = "processed"
	OutboxStatusFailed    OutboxStatus = "failed"
)

type OutboxMessage struct {
	ID         uint            `json:"id" gorm:"primaryKey"`
	Topic      string          `json:"topic" gorm:"not null"`
	Key        string          `json:"key"`
	Payload    json.RawMessage `json:"payload" gorm:"type:jsonb;not null"`
	Status     OutboxStatus    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	RetryCount int             `json:"retry_count" gorm:"default:0"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
