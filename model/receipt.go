package model

type Receipt struct {
	ID               int     `gorm:"column:id" json:"id"`
	ReceiptNumber    string  `gorm:"column:receipt_number" json:"receipt_number"`
	LoanID           int     `gorm:"column:loan_id" json:"loan_id"`
	ReceiptDate      string  `gorm:"column:receipt_date" json:"receipt_date"`
	TotalAmount      float64 `gorm:"column:total_amount" json:"total_amount"`
	CashierSessionID int     `gorm:"column:cashier_session_id" json:"cashier_session_id"`
	ReceiveBy        int     `gorm:"column:received_by" json:"received_by"`
	Notes            string  `gorm:"column:notes" json:"notes"`
}
