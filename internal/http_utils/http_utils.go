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

type GetModulesByIdsReq struct {
	ModulesIds []int `json:"modules_ids"`
}

type ResultForReq struct {
	Type     string               `json:"type"`
	CardsRes []entity.CardsResult `json:"cards_result,omitempty"`
}

type InsertModuleResultReq struct {
	ModuleId int          `json:"module_id"`
	Result   ResultForReq `json:"result"`
	Owner    int          `json:"owner"`
	Time     string       `json:"time,omitempty"`
}

type InsertCategoryModulesResultReq struct {
	CategoryId int                     `json:"category_id"`
	Modules    []InsertModuleResultReq `json:"modules_res"`
	Owner      int                     `json:"owner"`
	Time       string                  `json:"time"`
}

func GetModulesCreateReqFromJson(body []byte) (ModuleCreateReq, error) {
	var mod ModuleCreateReq
	err := json.Unmarshal(body, &mod)
	return mod, err
}
