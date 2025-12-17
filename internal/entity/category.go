package entity

type Category struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	OwnerId int      `json:"owner_id"`
	Modules []Module `json:"modules"`
}

type CategoryToCreate struct {
	Name    string `json:"name"`
	OwnerId int    `json:"-"`
	Modules []int  `json:"modules_ids"`
}
