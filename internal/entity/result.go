package entity

import "time"

type CardsResult struct {
	CardId int    `json:"card_id"`
	Result string `json:"result"`
}

type Result struct {
	Id       int           `json:"result_id"`
	Owner    int           `json:"owner"`
	Type     string        `json:"type"`
	Time     time.Time     `json:"time"`
	CardsRes []CardsResult `json:"cards_result,omitempty"`
}

type ModuleResult struct {
	ModuleId int    `json:"module_id"`
	Result   Result `json:"result,inline"`
}

type CategoryModulesResult struct {
	CategoryResultId int            `json:"category_result_id"`
	CategoryId       int            `json:"category_id"`
	Modules          []ModuleResult `json:"modules_res"`
}
