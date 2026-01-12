package model

type Loan struct {
	ID                 int     `json:"id"`
	ClientID           int     `json:"client_id"`
	CoID               int     `json:"co_id"`
	LoanProductID      int     `json:"loan_product_id"`
	LoanAmount         float32 `json:"loan_amount"`
	InterestRate       float32 `json:"interest_rate"`
	ProcessFee         float32 `json:"process_fee"`
	ApproveDate        *string `json:"approve_date"`
	LoanStartDate      *string `json:"loan_start_date"`
	LoanEndDate        *string `json:"loan_end_date"`
	DisbursedDate      string  `json:"disbursed_date"`
	DisbursedBy        int     `json:"disbursed_by"`
	DailyPaymentAmount float32 `json:"daily_payment_amount"`
	Purpose            string  `json:"purpose"`
	Duration           int     `json:"duration"`
	Status             Status  `json:"status"`
	DocumentTypeID     int     `json:"document_type_id"`
	CheckByID          int     `json:"check_by_id"`
	ApprovedByID       int     `json:"approved_by_id" gorm:"column:approved_by_id"`
	ClosedDate         *string `json:"closed_date"`
	ClosedReason       string  `json:"closed_reason"`
}
