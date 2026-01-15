package model

type ChartAccount struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	AccountTypeID int    `json:"account_type_id"`
	IsActive      bool   `json:"is_active" gorm:"column:is_active"`
}
