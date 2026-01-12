package request

type LoanRequest struct {
	ClientID       int                    `json:"client_id"`
	LoanProductID  int                    `json:"loan_product_id"`
	LoanAmount     float32                `json:"loan_amount"`
	Purpose        string                 `json:"purpose"`
	DocumentTypeID int                    `json:"document_type_id"`
	CheckByID      int                    `json:"check_by_id"`
	ApprovedByID   int                    `json:"approve_by_id"`
	Guarantors     []LoanGuarantorRequest `json:"guarantor_id"`
}
