package response

type Schedule struct {
	ID             int     `json:"id"`
	ScheduleNumber int     `json:"schedule_number"`
	PaymentDate    string  `json:"payment_date"`
	DueAmount      float64 `json:"due_amount"`
	Penalty        float64 `json:"penalty"`
	Total          float64 `json:"total"`
	PaidAmount     float64 `json:"paid_amount"`
	Totalowe       float64 `json:"total_owe"`
	Status         string  `json:"staus"`
}

type FullPaymentResponse struct {
	LoanID       int        `json:"loan_id"`
	ClientID     int        `json:"client_id"`
	ClientName   string     `json:"client_name"`
	ClientGender string     `json:"client_gender"`
	ClientPhone  string     `json:"client_phone"`
	LoanAmount   float64    `json:"loan_amount"`
	ProcessFee   float64    `json:"process_fee"`
	ApproveDate  *string    `json:"approve_date"`
	Purpose      string     `json:"purpose"`
	Duration     int        `json:"duration"`
	CoID         int        `json:"co_id"`
	CoName       string     `json:"co_name"`
	CoPhone      string     `json:"co_phone"`
	Schedule     []Schedule `json:"schedule" gorm:"-"`
}
