package entity

type User struct {
	Id           int        `json:"id"`
	Login        string     `json:"login,omitempty"`
	Name         string     `json:"name"`
	PasswordHash string     `json:"-"`
	Modules      []Module   `json:"modules"`
	Categories   []Category `json:"categories"`
}
