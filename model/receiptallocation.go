package model

type ReceiptAllocation struct {
	ID              int     `gorm:"column:id" json:"id"`
	ReceiptID       int     `gorm:"column:receipt_id" json:"receipt_id"`
	ScheduleID      int     `gorm:"column:schedule_id" json:"schedule_id"`
	PrincipalAmount float64 `gorm:"column:principal_amount" json:"principal_amount"`
	InterestAmount  float64 `gorm:"column:interest_amount" json:"interest_amount"`
	PenaltyAmount   float64 `gorm:"column:penalty_amount" json:"penalty_amount"`
	AllocationDate  string  `gorm:"column:allocation_date" json:"allocation_date"`
}
