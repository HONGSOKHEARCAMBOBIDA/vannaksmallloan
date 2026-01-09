package request

type AuthRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `jsos:"password" binding:"required"`
}
