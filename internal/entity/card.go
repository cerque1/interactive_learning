package entity

type TextWithLang struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type Card struct {
	Id           int          `json:"id"`
	ParentModule int          `json:"parent_module"`
	Term         TextWithLang `json:"term"`
	Definition   TextWithLang `json:"definition"`
}
