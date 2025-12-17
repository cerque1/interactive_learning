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
	GetUserInfoById(userId int) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
	InsertUser(user entity.User) error
}

type CardRepo interface {
	GetCardsByModule(moduleId int) ([]entity.Card, error)
	GetCardById(cardId int) (entity.Card, error)
	GetLastInsertedCardId() (int, error)
	InsertCard(card entity.Card) error
	DeleteCard(cardId int) error
}

type ModuleRepo interface {
	GetModulesByUser(userId int) ([]entity.Module, error)
	GetModuleById(moduleId int) (entity.Module, error)
	GetLastInsertedModuleId() (int, error)
	InsertModule(module entity.ModuleToCreate) error
	DeleteModule(moduleId int) error
}

type CategoryRepo interface {
	GetCategoriesToUser(userId int) ([]entity.Category, error)
	GetCategoryById(id int) (entity.Category, error)
	GetLastInsertedCategoryId() (int, error)
	InsertCategory(category entity.CategoryToCreate) error
	DeleteCategory(id int) error
}

type CategoryModulesRepo interface {
	GetModulesToCategory(categoryId int) ([]entity.Module, error)
	InsertModuleToCategory(categoryId, moduleId int) error
	DeleteModuleFromCategory(categoryId, moduleId int) error
	DeleteAllModulesFromCategory(categoryId int) error
}
