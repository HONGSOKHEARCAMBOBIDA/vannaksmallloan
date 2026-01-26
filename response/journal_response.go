package response

type JournalResponse struct {
	ID               int     `json:"id"`
	TransactionDate  string  `json:"transaction_date"`
	ChartAccountID   int     `json:"chart_account_id"`
	ChartAccountCode string  `json:"chart_account_code"`
	ChartAccountName string  `json:"chart_account_name"`
	DebitAmount      float64 `json:"debit_amount"`
	CreditAmount     float64 `json:"credit_amount"`
	Description      string  `json:"description"`
	ReferenceID      int     `json:"reference_id"`
	ReferenceCode    string  `json:"reference_code"`
	CreatedBy        int     `json:"created_by"`
	CreatedByName    string  `json:"created_by_name"`
}

type BalanceSheetTotals struct {
	TotalAssets            float64 `json:"total_assets"`
	TotalLiabilities       float64 `json:"total_liabilities"`
	TotalEquity            float64 `json:"total_equity"`
	TotalLiabilitiesEquity float64 `json:"total_liabilities_equity"`
	IsBalanced             bool    `json:"is_balanced"`
	Difference             float64 `json:"difference"`
}

type AccountBalance struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
}

type BalanceSheetSection struct {
	Title    string           `json:"title"`
	Accounts []AccountBalance `json:"accounts"`
	Total    float64          `json:"total"`
}

type BalanceSheetResponse struct {
	ReportTitle string              `json:"report_title"`
	ReportDate  string              `json:"report_date"`
	Assets      BalanceSheetSection `json:"assets"`
	Liabilities BalanceSheetSection `json:"liabilities"`
	Equity      BalanceSheetSection `json:"equity"`
	Totals      BalanceSheetTotals  `json:"totals"`
	Message     string              `json:"message"`
}
