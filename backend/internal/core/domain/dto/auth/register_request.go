package dto

type RegisterAccount struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}