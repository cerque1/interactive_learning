package entity

const (
	PrivateModule = 1
	PublicModule  = 0
)

type Module struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Cards   []Card `json:"cards,omitempty"`
	OwnerId int    `json:"owner_id"`
	Type    int    `json:"type"`
}

type PopularModule struct {
	Mod   Module `json:"module"`
	Count int    `json:"count"`
}

type ModuleToCreate struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Cards   []CardToAdd `json:"cards"`
	OwnerId int         `json:"owner_id"`
	Type    int         `json:"type"`
}
