package response

type RoleResponse struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `json:"name"`
	DisPlayName string `json:"display_name" gorm:"column:display_name"`
	IsActive    bool   `json:"is_active"`
}
