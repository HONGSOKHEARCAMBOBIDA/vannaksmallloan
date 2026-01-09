package request

type ClientRequestCreate struct {
	Name          string  `form:"name"`
	Gender        int     `form:"gender"`
	MaritatStatus string  `form:"marital_status"`
	DateofBirth   string  `form:"date_of_birth"`
	Occupation    string  `form:"occupation"`
	IdCardNumber  string  `form:"id_card_number"`
	Phone         string  `form:"phone"`
	VillageID     int     `form:"village_id"`
	Latitude      float64 `form:"latitude"`
	Longitude     float64 `form:"longitude"`
	// ImagePath     string  `form:"image_path"`
	Note string `form:"notes"`
}

type ClientRequestUpdate struct {
	Name          string  `form:"name"`
	Gender        int     `form:"gender"`
	MaritatStatus string  `form:"marital_status"`
	DateofBirth   string  `form:"date_of_birth"`
	Occupation    string  `form:"occupation"`
	IdCardNumber  string  `form:"id_card_number"`
	Phone         string  `form:"phone"`
	VillageID     int     `form:"village_id"`
	Latitude      float64 `form:"latitude"`
	Longitude     float64 `form:"longitude"`
	// ImagePath     string  `form:"image_path"`
	Note string `form:"notes"`
}
