package entity

const (
	PrivateModule = 0
	PublicModule  = 1
)

type Module struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Cards   []Card `json:"cards,omitempty"`
	OwnerId int    `json:"owner_id"`
	Type    int    `json:"type"`
}

type ModuleToCreate struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Cards   []CardToAdd `json:"cards"`
	OwnerId int         `json:"owner_id"`
	Type    int         `json:"type"`
}
