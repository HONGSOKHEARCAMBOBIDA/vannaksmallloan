package model

import "time"

type CashierSession struct {
	ID             int           `json:"id"`
	SessionNumber  string        `json:"session_number"`
	UserID         int           `json:"user_id"`
	SessionDate    string        `json:"session_date"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        *time.Time    `json:"end_time"`
	OpeningBalance float64       `json:"opening_balance"`
	ClosingBalance *float64      `json:"closing_balance"`
	TotalReceipts  *float64      `json:"total_receipts"`
	Difference     *float64      `json:"difference"`
	Status         StatusCashier `json:"status"`
	Notes          *string       `json:"notes"`
	VerifiedBy     *int          `json:"verified_by"`
	VerifiedAt     *string       `json:"verified_at"`
}
