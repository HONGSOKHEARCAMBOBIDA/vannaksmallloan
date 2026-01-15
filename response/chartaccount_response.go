package response

type ChartAccountResponse struct {
	ID              int    `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	AccountTypeID   int    `json:"account_type_id"`
	AccountTypeName string `json:"account_type_name"`
	Isactive        bool   `json:"is_active" gorm:"column:is_active"`
}
