package response

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username" gorm:"column:username"`
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Isactive bool   `json:"is_active" gorm:"column:is_active"`
}
