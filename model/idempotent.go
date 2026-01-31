package model

import "time"

type IdempotentRequest struct {
	ID        uint
	RequestID string `gorm:"uniqueIndex"`
	Endpoint  string
	CreatedAt time.Time
}
