package response

type CollectfromgoodloanResponse struct {
	ID             int     `json:"id"`
	ClientID       int     `json:"client_id"`
	ClientName     string  `json:"client_name"`
	UserID         int     `json:"user_id"`
	UserName       string  `json:"user_name"`
	VillageID      int     `json:"village_id"`
	VillageName    string  `json:"village_name"`
	TotalCollect   float64 `json:"total_collect" gorm:"column:total_collect"`
	TotalPenalty   float64 `json:"total_penalty" gorm:"column:total_penalty"`
	PenaltyDay     int     `json:"penalty_day" gorm:"column:penalty_day"`
	LumpSumpayment float64 `json:"lump_sum_payment" gorm:"column:lump_sum_payment"`
}

type ReceiptResponse struct {
	ID              int     `json:"id"`
	LoanID          int     `json:"loan_id"`
	ClientName      string  `json:"client_name"`
	CoName          string  `json:"co_name"`
	ReceiptDate     string  `json:"receipt_date"`
	TotalAmount     float64 `json:"total_amount"`
	Notes           string  `json:"notes"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestAmount  float64 `json:"interest_amount"`
	PenaltyAmount   float64 `json:"penalty_amount"`
	ReceiveByName   string  `json:"receive_by_name"`
}
