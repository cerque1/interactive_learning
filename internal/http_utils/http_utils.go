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

func GetModulesCreateReqFromJson(body []byte) (ModuleCreateReq, error) {
	var mod ModuleCreateReq
	err := json.Unmarshal(body, &mod)
	return mod, err
}
