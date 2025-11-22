package entity

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	PasswordHash string `json:"-"`
}
