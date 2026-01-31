package request

import "time"

type IdempotentRequest struct {
	RequestID string `gorm:"uniqueIndex"`
	Endpoint  string
	CreatedAt time.Time
}
