package usecase

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/utils/tokengenerator"
)

type Tokens interface {
	AddTokenToUser(id int) tokengenerator.Token
	DeleteTokenToUser(id int) error
	IsValidToken(token tokengenerator.Token) (int, error)
}

type Users interface {
	GetUserByLogin(login string) (entity.User, error)
	GetUserInfoById(user_id int, is_full bool) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
	InsertUser(user entity.User) (int, error)
}

type Cards interface {
	GetCardsByModule(module_id int) ([]entity.Card, error)
	GetCardById(card_id int) (entity.Card, error)
	InsertCard(card entity.Card) (int, error)
	InsertCards(cards []entity.Card) ([]int, error)
	DeleteCard(card_id int) error
}

type Modules interface {
	GetModulesByUser(user_id int) ([]entity.Module, error)
	GetModuleById(module_id int) (entity.Module, error)
	GetModulesWithCardsByUser(user_id int) ([]entity.Module, error)
	InsertModule(module entity.Module) (int, []int, error)
	DeleteModule(module_id int) error
}

type Categories interface {
	GetCategoriesToUser(user_id int, is_full bool) ([]entity.Category, error)
	GetCategoryById(id int) (entity.Category, error)
	InsertCategory(category entity.Category) (int, error)
	DeleteCategory(id int) error
}

type CategoryModules interface {
	GetModulesToCategory(category_id int, is_full bool) ([]entity.Module, error)
	InsertModuleToCategory(category_id, module_id int) error
	DeleteModuleFromCategory(category_id, module_id int) error
	DeleteAllModulesFromCategory(category_id int) error
}
