package model

type PaymentSchedule struct {
	ID              int            `gorm:"column:id" json:"id"`
	LoanID          int            `gorm:"column:loan_id" json:"loan_id"`
	ScheduleNumber  int            `gorm:"column:schedule_number" json:"schedule_number"`
	PaymentDate     string         `gorm:"column:payment_date" json:"payment_date"`
	DueAmount       float64        `gorm:"column:due_amount" json:"due_amount"`
	PrincipalAmount float64        `gorm:"column:principal_amount" json:"principal_amount"`
	InterestAmount  float64        `gorm:"column:interest_amount" json:"interest_amount"`
	PaidDate        *string        `gorm:"column:paid_date" json:"paid_date"`
	PrincipalPaid   *float64       `gorm:"column:principal_paid" json:"principal_paid"`
	InterestPaid    *float64       `gorm:"column:interest_paid" json:"interest_paid"`
	PenaltyAmount   float32        `gorm:"column:penalty_amount" json:"penalty_amount"`
	PenaltyPaid     *float32       `gorm:"column:penalty_paid" json:"penalty_paid"`
	PaidAmount      *float64       `gorm:"column:paid_amount" json:"paid_amount"`
	DayLate         *int           `gorm:"column:days_late" json:"days_late"`
	Status          ScheduelStatus `gorm:"column:status" json:"status"`
}
