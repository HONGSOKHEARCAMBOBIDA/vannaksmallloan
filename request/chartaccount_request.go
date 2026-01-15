package request

type ChartAccountRequestCreate struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	AccountTypeID int    `json:"account_type_id"`
}

type ChartAccountRequestUpdate struct {
	Code          *string `json:"code"`
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	AccountTypeID *int    `json:"account_type_id"`
}
