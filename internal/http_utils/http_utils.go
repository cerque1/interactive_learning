package httputils

import (
	"encoding/json"
	"interactive_learning/internal/entity"
)

type ModuleCreateReq struct {
	Name  string             `json:"name"`
	Type  int                `json:"type"`
	Cards []entity.CardToAdd `json:"cards"`
}

type AddModulesToCategoryReq struct {
	ModulesIds []int `json:"modules_ids"`
}

type RenameReq struct {
	NewName string `json:"new_name"`
}

func GetModulesCreateReqFromJson(body []byte) (ModuleCreateReq, error) {
	var mod ModuleCreateReq
	err := json.Unmarshal(body, &mod)
	return mod, err
}
