package model

type LoanProduct struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	InterestRate     float32 `json:"interest_rate"`
	ProcessFeeRate   float32 `json:"process_fee_rate"`
	TermDay          int     `json:"term_days" gorm:"column:term_days"`
	PaymentFrequency string  `json:"payment_frequency" gorm:"column:payment_frequency"`
	SkipWeeken       bool    `json:"skip_weekend" gorm:"column:skip_weekend"`
	LatePenaltyFixed float64 `json:"late_penalty_fixed" gorm:"column:late_penalty_fixed"`
	GracePeriodDay   int     `json:"grace_period_days" gorm:"column:grace_period_days"`
	Isactive         bool    `json:"is_active" gorm:"column:is_active"`
}
