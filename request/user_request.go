package request

type UserRequestCreate struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
	Phone    string `json:"phone"`
}
type UserRequestUpdate struct {
	Name     *string `json:"name,omitempty"`
	UserName *string `json:"username,omitempty"`
	RoleId   *int    `json:"role_id,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}
type UserRequestResetPassword struct {
	NewPassword string `json:"new_password" binding:"required"`
}
