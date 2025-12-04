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
	GetUserInfoById(userId int, isFull bool) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
	InsertUser(user entity.User) (int, error)
}

type Cards interface {
	GetCardsByModule(moduleId int) ([]entity.Card, error)
	GetCardById(cardId int) (entity.Card, error)
	InsertCard(card entity.Card) (int, error)
	InsertCards(cards []entity.Card) ([]int, error)
	DeleteCard(cardId int) error
}

type Modules interface {
	GetModulesByUser(userId int) ([]entity.Module, error)
	GetModuleById(moduleId int) (entity.Module, error)
	GetModulesWithCardsByUser(userId int) ([]entity.Module, error)
	InsertModule(module entity.Module) (int, []int, error)
	DeleteModule(moduleId int) error
}

type Categories interface {
	GetCategoriesToUser(userId int, isFull bool) ([]entity.Category, error)
	GetCategoryById(id int) (entity.Category, error)
	InsertCategory(category entity.Category) (int, error)
	DeleteCategory(id int) error
}

type CategoryModules interface {
	GetModulesToCategory(categoryId int, isFull bool) ([]entity.Module, error)
	InsertModuleToCategory(categoryId, moduleId int) error
	DeleteModuleFromCategory(categoryId, moduleId int) error
	DeleteAllModulesFromCategory(categoryId int) error
}
