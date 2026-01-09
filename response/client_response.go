package response

type ClientResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Gender        int     `json:"gender"`
	MaritatStatus string  `json:"marital_status"`
	DateofBirth   string  `json:"date_of_birth"`
	Occupation    string  `json:"occupation"`
	IdCardNumber  string  `json:"id_card_number"`
	Phone         string  `json:"phone"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	ImagePath     string  `json:"image_path"`
	Note          string  `json:"notes"`
	Isactive      bool    `json:"is_active"`
	CreateBy      int     `json:"created_by"`
	CreateByName  string  `json:"create_by_name"`
	ProvinceID    int     `json:"province_id"`
	ProvinceName  string  `json:"province_name"`
	DistrictID    int     `json:"district_id"`
	DistrictName  string  `json:"district_name"`
	CommunceID    int     `json:"communce_id"`
	CommunceName  string  `json:"communce_name"`
	VillageID     int     `json:"village_id"`
	VillageName   string  `json:"village_name"`
}
