package models

type TokenData struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Token    string `json:"token"`
}
