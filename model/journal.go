package model

type Journal struct {
	ID              int     `json:"id"`
	TransactionDate string  `json:"transaction_date"`
	ChartAccountID  int     `json:"chart_account_id"`
	DebitAmount     float64 `json:"debit_amount"`
	CreditAmount    float64 `json:"credit_amount"`
	Description     string  `json:"description"`
	ReferenceID     int     `json:"reference_id"`
	ReferenceCode   string  `json:"reference_code"`
	CreatedBy       int     `json:"created_by"`
}
