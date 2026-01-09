package response

type AuthResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `josn:"token"`
}
