package usecase

import (
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
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
	// GetCardOwnerId(cardId int) (int, error)
	InsertCard(card entity.Card) (int, error)
	InsertCards(cards entity.CardsToAdd) ([]int, error)
	UpdateCard(userId int, card entity.Card) error
	DeleteCard(userId int, cardId int) error
	// DeleteCardsToParentModule(moduleId int) error
}

type Modules interface {
	GetModulesByUser(userId int) ([]entity.Module, error)
	GetModuleById(moduleId int) (entity.Module, error)
	GetModulesByIds(modulesIds []int, isFull bool) ([]entity.Module, error)
	GetModulesWithCardsByUser(userId int) ([]entity.Module, error)
	GetModuleOwnerId(moduleId int) (int, error)
	InsertModule(module entity.ModuleToCreate) (int, []int, error)
	RenameModule(userId, moduleId int, newName string) error
	DeleteModule(userId int, moduleId int) error
}

type Categories interface {
	GetCategoriesToUser(userId int, isFull bool) ([]entity.Category, error)
	GetCategoryById(id int) (entity.Category, error)
	InsertCategory(category entity.CategoryToCreate) (int, error)
	RenameCategory(userId, categoryId int, newName string) error
	DeleteCategory(userId, categoryId int) error
}

type CategoryModules interface {
	GetModulesToCategory(categoryId int, isFull bool) ([]entity.Module, error)
	InsertModulesToCategory(userId, categoryId int, modulesIds []int) error
	DeleteModuleFromCategory(userId, categoryId, moduleId int) error
	// DeleteAllModulesFromCategory(userId, categoryId int) error
	// DeleteModuleFromCategories(moduleId int) error
}

type Results interface {
	GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error)
	GetCardsResultById(resultId int) ([]entity.CardsResult, error)
	GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error)
	GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error)
	GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error)
	InsertModuleResult(result httputils.InsertModuleResultReq) (int, error)
	InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) ([]int, error)
	DeleteModuleResult(resultId int) error
	DeleteCategoryResultById(categoryResultId int) error
	// DeleteResultByModuleId(moduleId int) error
	// DeleteResultByCategoryId(categoryId int) error
}
