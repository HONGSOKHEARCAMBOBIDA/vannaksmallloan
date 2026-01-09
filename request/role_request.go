package request

type RoleRequestUpdate struct {
	Name        *string `json:"name"`
	DisPlayName *string `json:"display_name" gorm:"column:display_name"`
}
type RoleRequestCreate struct {
	Name        string `json:"name"`
	DisPlayName string `json:"display_name" gorm:"column:display_name"`
	IsActive    bool   `json:"is_active"`
}
