package response

type CollectfromgoodloanResponse struct {
	ID           int     `json:"id"`
	ClientID     int     `json:"client_id"`
	ClientName   string  `json:"client_name"`
	UserID       int     `json:"user_id"`
	UserName     string  `json:"user_name"`
	VillageID    int     `json:"village_id"`
	VillageName  string  `json:"village_name"`
	TotalCollect float64 `json:"total_collect" gorm:"column:total_collect"`
	TotalPenalty float64 `json:"total_penalty" gorm:"column:total_penalty"`
}
