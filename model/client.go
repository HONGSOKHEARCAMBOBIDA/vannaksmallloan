package model

type Client struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Gender        Gender  `json:"gender"`
	MaritatStatus string  `json:"marital_status" gorm:"column:marital_status"`
	DateOfBirth   string  `json:"date_of_birth"`
	Occupation    string  `json:"occupation"`
	IdCardNumber  string  `json:"id_card_number"`
	Phone         string  `json:"phone"`
	VillageID     int     `json:"village_id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	ImagePath     string  `json:"image_path"`
	Note          string  `json:"notes" gorm:"column:notes"`
	IsActive      bool    `json:"is_active"`
	CreateBy      int     `json:"created_by" gorm:"column:created_by"`
}
