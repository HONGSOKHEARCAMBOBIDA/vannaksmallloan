package response

type LoanResponse struct {
	ID                  int     `json:"id"`
	ClientID            int     `json:"client_id"`
	ClientName          string  `json:"client_name"`
	ClientGender        int     `json:"client_gender"`
	ClientMaritalStatus string  `json:"client_marital_status"`
	ClientDoB           string  `json:"client_date_of_birth" gorm:"column:client_date_of_birth"`
	ClientOccupation    string  `json:"client_occupation"`
	ClientPhone         string  `json:"client_phone"`
	ProvinceID          int     `json:"province_id"`
	ProvinceName        string  `json:"province_name"`
	DistrictID          int     `json:"district_id"`
	DistrictName        string  `json:"district_name"`
	CommunceID          int     `json:"communce_id"`
	CommunceName        string  `json:"communce_name"`
	VillageID           int     `json:"village_id"`
	VillageName         string  `json:"village_name"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	CoID                int     `json:"co_id"`
	CoName              string  `json:"co_name"`
	LoanProductID       int     `json:"loan_product_id"`
	LoanProductName     string  `json:"loan_product_name"`
	LoanAmount          float64 `json:"loan_amount"`
	InterestRate        float32 `json:"interest_rate"`
	ProcessFee          float32 `json:"process_fee"`
	ApproveDate         string  `json:"approve_date"`
	LoanStartDate       string  `json:"loan_start_date"`
	LoanEndDate         string  `json:"loan_end_date"`
	DisbursedDate       string  `json:"disbursed_date"`
	DisbursedBy         int     `json:"disbursed_by"`
	DisbursedByName     string  `json:"disburse_by_name" gorm:"column:disburse_by_name"`
	DailyPaymentAmount  float32 `json:"daily_payment_amount"`
	Purpose             string  `json:"purpose"`
	Duration            string  `json:"duration"`
	Status              string  `json:"status"`
	DocumentTypeID      int     `json:"document_type_id"`
	DocumentTypeName    string  `json:"document_type_name"`
	CheckByID           int     `json:"check_by_id"`
	CheckByName         string  `json:"check_by_name"`
	ApprovedByID        int     `json:"approve_by_id" gorm:"column:approve_by_id"`
	ApproveByName       string  `json:"approve_by_name"`
	ClosedDate          string  `json:"close_date"`
	ClosedReason        string  `json:"close_reason"`
}
