package response

import "time"

type CashierSessionResponse struct {
	ID             int       `json:"id"`
	SessionNumber  string    `json:"session_number"`
	UserID         int       `json:"user_id"`
	UserName       string    `json:"user_name"`
	SessionDate    string    `json:"session_date"`
	StartDate      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	OpeningBalance float64   `json:"opening_balance"`
	ClosingBalance float64   `json:"closing_balance"`
	TotalReceipt   float64   `json:"total_receipts"`
	Difference     float64   `json:"difference"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	VerifiedBy     int       `json:"verified_by"`
	VerifiedByName string    `json:"verified_by_name"`
	VerifiedAt     string    `json:"verified_at"`
}
