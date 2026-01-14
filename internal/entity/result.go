package entity

import "time"

type CardsResult struct {
	CardId int    `json:"card_id"`
	Result string `json:"result"`
}

type Result struct {
	Id       int           `json:"result_id"`
	Type     string        `json:"type"`
	CardsRes []CardsResult `json:"cards_result,omitempty"`
}

type ModuleResult struct {
	ModuleId int        `json:"module_id"`
	Owner    int        `json:"owner,omitempty"`
	Time     *time.Time `json:"time,omitempty"`
	Result   Result     `json:"result"`
}

type CategoryModulesResult struct {
	CategoryResultId int            `json:"category_result_id"`
	CategoryId       int            `json:"category_id"`
	Owner            int            `json:"owner"`
	Time             time.Time      `json:"time"`
	Modules          []ModuleResult `json:"modules_res"`
}
