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

type CardToAdd struct {
	Id         int          `json:"id"`
	Term       TextWithLang `json:"term"`
	Definition TextWithLang `json:"definition"`
}

type CardsToAdd struct {
	Cards        []CardToAdd `json:"cards"`
	ParentModule int         `json:"parent_module"`
}
