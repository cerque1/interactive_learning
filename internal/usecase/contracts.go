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
	GetUsersWithSimilarName(name string, limit, offset int) ([]entity.User, error)
	GetUserByLogin(login string) (entity.User, error)
	GetUserInfoById(ownerId int, isFull bool, userId int) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
	InsertUser(user entity.User) (int, error)
}

type Cards interface {
	GetCardsByModule(moduleId, userId int) ([]entity.Card, error)
	GetCardById(cardId int) (entity.Card, error)
	InsertCard(card entity.Card) (int, error)
	InsertCards(cards entity.CardsToAdd) ([]int, error)
	UpdateCard(userId int, card entity.Card) error
	DeleteCard(userId int, cardId int) error
}

type Modules interface {
	GetModulesWithSimilarName(name string, limit, offset, userId int) ([]entity.Module, error)
	GetModulesByUser(ownerId int, withCards bool, userId int) ([]entity.Module, error)
	GetModuleById(moduleId, userId int) (entity.Module, error)
	GetModulesByIds(modulesIds []int, isFull bool, userId int) ([]entity.Module, error)
	GetModuleOwnerId(moduleId int) (int, error)
	InsertModule(module entity.ModuleToCreate) (int, []int, error)
	RenameModule(userId, moduleId int, newName string) error
	UpdateModuleType(moduleId, newType, userId int) error
	DeleteModule(userId int, moduleId int) error
}

type Categories interface {
	GetCategoriesWithSimilarName(name string, limit, offset, userId int) ([]entity.Category, error)
	GetCategoriesToUser(ownerId int, isFull bool, userId int) ([]entity.Category, error)
	GetCategoryById(id, userId int) (entity.Category, error)
	InsertCategory(category entity.CategoryToCreate) (int, error)
	RenameCategory(userId, categoryId int, newName string) error
	UpdateCategoryType(categoryId, newType, userId int) error
	DeleteCategory(userId, categoryId int) error
}

type CategoryModules interface {
	GetModulesToCategory(categoryId int, isFull bool, userId int) ([]entity.Module, error)
	InsertModulesToCategory(userId, categoryId int, modulesIds []int) error
	DeleteModuleFromCategory(userId, categoryId, moduleId int) error
}

type Results interface {
	GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error)
	GetModuleResultById(resultId int) (entity.ModuleResult, error)
	GetCardsResultById(resultId int) ([]entity.CardsResult, error)
	GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error)
	GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error)
	GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error)
	InsertModuleResult(result httputils.InsertModuleResultReq) (int, error)
	InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) (int, []int, error)
	DeleteModuleResult(resultId int) error
	DeleteCategoryResultById(categoryResultId int) error
}
