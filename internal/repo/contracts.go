package repo

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/utils/tokengenerator"
)

type TokenStorage interface {
	AddTokenToUser(id int) tokengenerator.Token
	DeleteTokenToUser(id int) error
	IsValidToken(token tokengenerator.Token) (int, error)
}

type UsersRepo interface {
	GetUserByLogin(login string) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
	InsertUser(user entity.User) error
}

type CardRepo interface {
	GetCardsByModule(module_id int) ([]entity.Card, error)
	GetCardById(card_id int) (entity.Card, error)
	GetLastInsertedCardId() (int, error)
	InsertCard(card entity.Card) error
	DeleteCard(card_id int) error
}

type ModuleRepo interface {
	GetModulesByUser(user_id int) ([]entity.Module, error)
	GetModuleById(module_id int) (entity.Module, error)
	GetLastInsertedModuleId() (int, error)
	InsertModule(module entity.Module) error
	DeleteModule(module_id int) error
}
