package request

type JournalRequestCreate struct {
	TransactionDate string  `json:"transaction_date"`
	DebitAccountID  int     `json:"debit_account_id"`
	CreditAccountID int     `json:"credit_account_id"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
}

type JournalRequestUpdate struct {
	TransactionDate *string  `json:"transaction_date"`
	ChartAccountID  *int     `json:"chart_account_id"`
	DebitAmount     *float64 `json:"debit_amount"`
	CreditAmount    *float64 `json:"credit_amount"`
	Description     *string  `json:"description"`
}
