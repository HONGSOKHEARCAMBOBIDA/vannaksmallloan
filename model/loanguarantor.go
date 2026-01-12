package model

type LoanGuarantor struct {
	ID           int    `json:"id"`
	LoanID       int    `json:"loan_id"`
	ClientID     int    `json:"client_id"`
	Relationship string `json:"relationship"`
	SignedDate   string `json:"signed_date"`
	Notes        string `json:"notes"`
}
