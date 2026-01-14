package repo

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/utils/tokengenerator"
	"time"
)

type TokenStorage interface {
	AddTokenToUser(id int) tokengenerator.Token
	DeleteTokenToUser(id int) error
	IsValidToken(token tokengenerator.Token) (int, error)
}

type UsersRepoRead interface {
	GetUserByLogin(login string) (entity.User, error)
	GetUserInfoById(userId int) (entity.User, error)
	IsContainsLogin(login string) (bool, error)
}

type UsersRepoWrite interface {
	InsertUser(user entity.User) error
}

type CardRepoRead interface {
	GetCardsByModule(moduleId int) ([]entity.Card, error)
	GetCardById(cardId int) (entity.Card, error)
	GetLastInsertedCardId() (int, error)
	GetParentModuleId(cardId int) (int, error)
}

type CardRepoWrite interface {
	InsertCard(card entity.Card) error
	UpdateCard(card entity.Card) error
	DeleteCard(cardId int) error
	DeleteCardsToParentModule(moduleId int) error
}

type ModuleRepoRead interface {
	GetModulesByUser(userId int) ([]entity.Module, error)
	GetModuleById(moduleId int) (entity.Module, error)
	GetLastInsertedModuleId() (int, error)
	GetModuleOwnerId(moduleId int) (int, error)
}

type ModuleRepoWrite interface {
	InsertModule(module entity.ModuleToCreate) error
	RenameModule(moduleId int, newName string) error
	DeleteModule(moduleId int) error
}

type CategoryRepoRead interface {
	GetCategoriesToUser(userId int) ([]entity.Category, error)
	GetCategoryById(id int) (entity.Category, error)
	GetLastInsertedCategoryId() (int, error)
	GetCategoryOwnerId(categoryId int) (int, error)
}

type CategoryRepoWrite interface {
	InsertCategory(category entity.CategoryToCreate) error
	RenameCategory(categoryId int, newName string) error
	DeleteCategory(categoryId int) error
}

type CategoryModulesRepoRead interface {
	GetModulesToCategory(categoryId int) ([]entity.Module, error)
}

type CategoryModulesRepoWrite interface {
	InsertModulesToCategory(categoryId, moduleId int) error
	DeleteModuleFromCategory(categoryId, moduleId int) error
	DeleteAllModulesFromCategory(categoryId int) error
	DeleteModuleFromCategories(moduleId int) error
}

type ResultsRepoRead interface {
	GetResultsByOwner(ownerId int) ([]entity.Result, error)
	GetResultById(id int) (entity.Result, error)
	GetLastInsertedResultId() (int, error)
}

type ResultsRepoWrite interface {
	InsertResult(result entity.Result) error
	DeleteResultById(id int) error
}

type CardsResultsRepoRead interface {
	GetCardsResultById(resultId int) ([]entity.CardsResult, error)
}

type CardsResultsRepoWrite interface {
	InsertCardResult(resultId, cardId int, result string) error
	DeleteCardResult(resultId, cardId int) error
	DeleteCardsToResult(resultId int) error
	DeleteResultsToCard(cardId int) error
}

type ModulesResultsRepoRead interface {
	GetModulesResultById(resultId int) (entity.ModuleResult, error)
	GetModulesResByOwner(ownerId int) ([]entity.ModuleResult, error)
	GetResultsToModuleOwner(moduleId, ownerId int) ([]entity.ModuleResult, error)
	GetResultsToModule(moduleId int) ([]entity.ModuleResult, error)
}

type ModulesResultsRepoWrite interface {
	InsertResultToModule(moduleId, resultId, ownerId int, time time.Time) error
	DeleteResultsToModule(moduleId int) error
	DeleteResultToModule(resultId int) error
}

type CategoryModulesResultsRepoRead interface {
	GetCategoriesResByOwner(ownerId int) ([]entity.CategoryModulesResult, error)
	GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error)
	GetResultsByCategoryOwner(categoryId, ownerId int) ([]entity.CategoryModulesResult, error)
	GetResultsByCategoryAndModule(categoryId, moduleId int) ([]int, error)
	GetLastInsertedResId() (int, error)
	GetResultsByModuleId(moduleId int) ([]int, error)
}

type CategoryModulesResultsRepoWrite interface {
	InsertCategoryModule(categoryResultId, categoryId, moduleId, result_id, ownerId int, time time.Time) error
	DeleteModulesFromCategories(moduleId int) error
	DeleteModulesFromCategory(categoryId, moduleId int) error
	DeleteAllToCategory(categoryId int) error
	DeleteResultById(categoryResultId int) error
}
