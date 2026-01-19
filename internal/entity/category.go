package entity

const (
	PublicCategory  = 0
	PrivateCategory = 1
)

type Category struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	OwnerId int      `json:"owner_id"`
	Modules []Module `json:"modules,omitempty"`
	Type    int      `json:"type"` // если >= 1 то приват, значение больше храниться в бд для отслеживания сколько приватных модулей
}

type PopularCategory struct {
	Cat   Category `json:"category"`
	Count int      `json:"users_count"`
}

type CategoryToCreate struct {
	Name    string `json:"name"`
	OwnerId int    `json:"-"`
	Modules []int  `json:"modules_ids"`
	Type    int    `json:"type"`
}
