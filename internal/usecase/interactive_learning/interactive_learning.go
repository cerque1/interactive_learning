package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/repo/persistent"
	"interactive_learning/internal/utils/tokengenerator"
	"slices"
	"sync"
	"time"
)

type UseCase struct {
	userRepo                   repo.UsersRepo
	tokenStorage               repo.TokenStorage
	cardRepo                   repo.CardRepo
	moduleRepo                 repo.ModuleRepo
	categoryRepo               repo.CategoryRepo
	categoryModulesRepo        repo.CategoryModulesRepo
	resultRepo                 repo.ResultsRepo
	cardsResultsRepo           repo.CardsResultsRepo
	modulesResultsRepo         repo.ModulesResultsRepo
	categoryModulesResultsRepo repo.CategoryModulesResultsRepo

	usersMutex                  sync.Mutex
	cardMutex                   sync.Mutex
	moduleMutex                 sync.Mutex
	categoryMutex               sync.Mutex
	categoryModulesMutex        sync.Mutex
	resultsMutex                sync.Mutex
	cardsResultsMutex           sync.Mutex
	modulesResultsMutex         sync.Mutex
	categoryModulesResultsMutex sync.Mutex
}

func New(userRepo repo.UsersRepo, cardRepo repo.CardRepo, moduleRepo repo.ModuleRepo, categoryRepo repo.CategoryRepo, categoryModulesRepo repo.CategoryModulesRepo) *UseCase {
	return &UseCase{userRepo: userRepo,
		tokenStorage:        persistent.NewTokenStorage(),
		cardRepo:            cardRepo,
		moduleRepo:          moduleRepo,
		categoryRepo:        categoryRepo,
		categoryModulesRepo: categoryModulesRepo,
	}
}

func (u *UseCase) AddTokenToUser(id int) tokengenerator.Token {
	return u.tokenStorage.AddTokenToUser(id)
}

func (u *UseCase) DeleteTokenToUser(id int) error {
	return u.tokenStorage.DeleteTokenToUser(id)
}

func (u *UseCase) IsValidToken(token tokengenerator.Token) (int, error) {
	return u.tokenStorage.IsValidToken(token)
}

func (u *UseCase) GetUserByLogin(login string) (entity.User, error) {
	return u.userRepo.GetUserByLogin(login)
}

func (u *UseCase) GetUserInfoById(userId int, isFull bool) (entity.User, error) {
	user, err := u.userRepo.GetUserInfoById(userId)
	if err != nil {
		return entity.User{}, err
	}

	if !isFull {
		return user, nil
	}

	modules, err := u.GetModulesWithCardsByUser(userId)
	if err != nil {
		return entity.User{}, err
	}
	user.Modules = modules

	categories, err := u.GetCategoriesToUser(userId, true)
	if err != nil {
		return entity.User{}, err
	}
	user.Categories = categories

	return user, nil
}

func (u *UseCase) GetCardOwnerId(cardId int) (int, error) {
	parentModule, err := u.cardRepo.GetParentModuleId(cardId)
	if err != nil {
		return -1, err
	}
	moduleOwner, err := u.moduleRepo.GetModuleOwnerId(parentModule)
	if err != nil {
		return -1, err
	}
	return moduleOwner, nil
}

func (u *UseCase) IsContainsLogin(login string) (bool, error) {
	u.usersMutex.Lock()
	defer u.usersMutex.Unlock()

	return u.userRepo.IsContainsLogin(login)
}

func (u *UseCase) InsertUser(user entity.User) (int, error) {
	u.usersMutex.Lock()
	defer u.usersMutex.Unlock()

	err := u.userRepo.InsertUser(user)
	if err != nil {
		return -1, err
	}
	newUser, err := u.userRepo.GetUserByLogin(user.Login)
	if err != nil {
		return -1, err
	}
	return newUser.Id, nil
}

func (u *UseCase) GetCardById(cardId int) (entity.Card, error) {
	return u.cardRepo.GetCardById(cardId)
}

func (u *UseCase) GetCardsByModule(moduleId int) ([]entity.Card, error) {
	return u.cardRepo.GetCardsByModule(moduleId)
}

func (u *UseCase) InsertCard(card entity.Card) (int, error) {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	err := u.cardRepo.InsertCard(card)
	if err != nil {
		return -1, err
	}
	return u.cardRepo.GetLastInsertedCardId()
}

func (u *UseCase) InsertCards(cards entity.CardsToAdd) ([]int, error) {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	var ids []int
	var err error
	var curId int

	for _, card := range cards.Cards {
		err = u.cardRepo.InsertCard(entity.Card{ParentModule: cards.ParentModule, Term: card.Term, Definition: card.Definition})
		if err != nil {
			return []int{}, err
		}
		curId, err = u.cardRepo.GetLastInsertedCardId()
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, curId)
	}
	return ids, nil
}

func (u *UseCase) UpdateCard(userId int, card entity.Card) error {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.GetCardOwnerId(card.Id)
	if err != nil {
		return errors.New("bad card id")
	} else if ownerId != userId {
		return errors.New("unaccessable card")
	}

	return u.cardRepo.UpdateCard(card)
}

func (u *UseCase) DeleteCard(userId int, cardId int) error {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.GetCardOwnerId(cardId)
	if err != nil {
		return errors.New("bad card id")
	}

	if ownerId != userId {
		return errors.New("unaccessable card")
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err = u.cardsResultsRepo.DeleteResultsToCard(cardId)
	if err != nil {
		return err
	}

	return u.cardRepo.DeleteCard(cardId)
}

func (u *UseCase) DeleteCardsToParentModule(moduleId int) error {
	module, err := u.GetModuleById(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, card := range module.Cards {
		err = u.cardsResultsRepo.DeleteResultsToCard(card.Id)
		if err != nil {
			return err
		}
	}

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	return u.cardRepo.DeleteCardsToParentModule(moduleId)
}

func (u *UseCase) GetModulesByUser(userId int) ([]entity.Module, error) {
	return u.moduleRepo.GetModulesByUser(userId)
}

func (u *UseCase) GetModulesWithCardsByUser(userId int) ([]entity.Module, error) {
	modules, err := u.moduleRepo.GetModulesByUser(userId)
	if err != nil {
		return []entity.Module{}, err
	}

	for i := range modules {
		cards, err := u.cardRepo.GetCardsByModule(modules[i].Id)
		if err != nil {
			return []entity.Module{}, err
		}
		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) GetModuleById(moduleId int) (entity.Module, error) {
	module, err := u.moduleRepo.GetModuleById(moduleId)
	if err != nil {
		return entity.Module{}, err
	}
	cards, err := u.cardRepo.GetCardsByModule(moduleId)
	if err != nil {
		return entity.Module{}, err
	}
	module.Cards = cards
	return module, nil
}

func (u *UseCase) GetModulesByIds(modulesIds []int, isFull bool) ([]entity.Module, error) {
	modules := []entity.Module{}

	for _, moduleId := range modulesIds {
		module, err := u.moduleRepo.GetModuleById(moduleId)
		if err != nil {
			return []entity.Module{}, err
		}

		if isFull {
			cards, err := u.cardRepo.GetCardsByModule(moduleId)
			if err != nil {
				return []entity.Module{}, err
			}
			module.Cards = cards
		}

		modules = append(modules, module)
	}

	return modules, nil
}

func (u *UseCase) GetModuleOwnerId(moduleId int) (int, error) {
	return u.moduleRepo.GetModuleOwnerId(moduleId)
}

func (u *UseCase) InsertModule(module entity.ModuleToCreate) (int, []int, error) {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	err := u.moduleRepo.InsertModule(module)
	if err != nil {
		return -1, []int{}, err
	}
	insertIds, err := u.InsertCards(entity.CardsToAdd{Cards: module.Cards, ParentModule: module.Id})
	if err != nil {
		return -1, []int{}, err
	}
	id, err := u.moduleRepo.GetLastInsertedModuleId()
	if err != nil {
		return -1, []int{}, err
	}
	return id, insertIds, nil
}

func (u *UseCase) RenameModule(userId, moduleId int, newName string) error {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	ownerId, err := u.GetModuleOwnerId(moduleId)
	if err != nil {
		return errors.New("bad module id")
	}
	if ownerId != userId {
		return errors.New("unaccessable module")
	}

	return u.moduleRepo.RenameModule(moduleId, newName)
}

func (u *UseCase) DeleteModule(userId int, moduleId int) error {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	ownerId, err := u.GetModuleOwnerId(moduleId)
	if err != nil {
		return errors.New("bad module id")
	}
	if ownerId != userId {
		return errors.New("unaccessable module")
	}

	err = u.DeleteCardsToParentModule(moduleId)
	if err != nil {
		return err
	}

	err = u.DeleteModuleFromCategories(moduleId)
	if err != nil {
		return err
	}

	err = u.DeleteResultByModuleId(moduleId)
	if err != nil {
		return err
	}

	err = u.DeleteModuleResFromCategories(moduleId)
	if err != nil {
		return err
	}

	return u.moduleRepo.DeleteModule(moduleId)
}

func (u *UseCase) GetCategoriesToUser(userId int, isFull bool) ([]entity.Category, error) {
	categories, err := u.categoryRepo.GetCategoriesToUser(userId)
	if err != nil {
		return []entity.Category{}, err
	}

	if !isFull {
		return categories, nil
	}

	for i := range categories {
		modules, err := u.GetModulesToCategory(categories[i].Id, true)
		if err != nil {
			return []entity.Category{}, err
		}

		categories[i].Modules = modules
	}

	return categories, nil
}

func (u *UseCase) GetCategoryById(id int) (entity.Category, error) {
	category, err := u.categoryRepo.GetCategoryById(id)
	if err != nil {
		return entity.Category{}, nil
	}
	modules, err := u.GetModulesToCategory(id, true)
	if err != nil {
		return entity.Category{}, err
	}
	category.Modules = modules

	return category, nil
}

func (u *UseCase) InsertCategory(category entity.CategoryToCreate) (int, error) {
	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	err := u.categoryRepo.InsertCategory(category)
	if err != nil {
		return -1, err
	}
	new_id, err := u.categoryRepo.GetLastInsertedCategoryId()
	if err != nil {
		return -1, err
	}
	if err = u.InsertModulesToCategory(category.OwnerId, new_id, category.Modules); err != nil {
		return -1, err
	}
	return new_id, nil
}

func (u *UseCase) IsCategoryOwner(userId, categoryId int) (bool, error) {
	ownerId, err := u.categoryRepo.GetCategoryOwnerId(categoryId)
	if err != nil {
		return false, errors.New("bad category id")
	}
	if ownerId != userId {
		return false, nil
	}
	return true, nil
}

func (u *UseCase) RenameCategory(userId, categoryId int, newName string) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.IsCategoryOwner(userId, categoryId)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	return u.categoryRepo.RenameCategory(categoryId, newName)
}

func (u *UseCase) DeleteCategory(userId int, id int) error {
	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	isOwner, err := u.IsCategoryOwner(userId, id)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	err = u.DeleteAllModulesFromCategory(userId, id)
	if err != nil {
		return err
	}

	err = u.DeleteResultByCategoryId(id)
	if err != nil {
		return err
	}

	return u.categoryRepo.DeleteCategory(id)
}

func (u *UseCase) GetModulesToCategory(categoryId int, isFull bool) ([]entity.Module, error) {
	modules, err := u.categoryModulesRepo.GetModulesToCategory(categoryId)
	if err != nil {
		return []entity.Module{}, err
	}

	if !isFull {
		return modules, nil
	}

	for i := range modules {
		cards, err := u.cardRepo.GetCardsByModule(modules[i].Id)
		if err != nil {
			return []entity.Module{}, err
		}

		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) InsertModulesToCategory(userId, categoryId int, modulesIds []int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	category, err := u.categoryRepo.GetCategoryById(categoryId)
	if err != nil {
		return errors.New("bad category id")
	}
	if category.OwnerId != userId {
		return errors.New("unavailable category")
	}

	for _, moduleId := range modulesIds {
		if idx := slices.IndexFunc(category.Modules, func(elt entity.Module) bool { return elt.Id == moduleId }); idx >= 0 {
			return errors.New("module is already exists")
		}
	}

	for _, moduleId := range modulesIds {
		err := u.categoryModulesRepo.InsertModulesToCategory(categoryId, moduleId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) DeleteModuleFromCategory(userId, categoryId, moduleId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.IsCategoryOwner(userId, categoryId)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	err = u.DeleteModuleResFromCategory(categoryId, moduleId)
	if err != nil {
		return err
	}
	return u.categoryModulesRepo.DeleteModuleFromCategory(categoryId, moduleId)
}

func (u *UseCase) DeleteAllModulesFromCategory(userId, categoryId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.IsCategoryOwner(userId, categoryId)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	return u.categoryModulesRepo.DeleteAllModulesFromCategory(categoryId)
}

func (u *UseCase) DeleteModuleFromCategories(moduleId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.DeleteModuleFromCategories(moduleId)
}

func (u *UseCase) GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error) {
	categoriesRes, err := u.categoryModulesResultsRepo.GetCategoriesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	modulesRes, err := u.modulesResultsRepo.GetModulesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	return categoriesRes, modulesRes, nil
}

func (u *UseCase) GetCardsResultById(resultId int) ([]entity.CardsResult, error) {
	return u.cardsResultsRepo.GetCardsResultById(resultId)
}

func (u *UseCase) GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error) {
	return u.modulesResultsRepo.GetResultsToModuleOwner(moduleId, userId)
}

func (u *UseCase) GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepo.GetResultsByCategoryOwner(categoryId, userId)
}

func (u *UseCase) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepo.GetCategoryResById(categoryResultsId)
}

func (u *UseCase) InsertModuleResult(result httputils.InsertModuleResultReq) (int, error) {
	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	time, err := time.Parse(time.RFC3339, result.Result.Time)
	if err != nil {
		return -1, err
	}

	err = u.resultRepo.InsertResult(entity.Result{Owner: result.Result.Owner,
		Type: result.Result.Type,
		Time: time})
	if err != nil {
		return -1, err
	}

	insertedResId, err := u.resultRepo.GetLastInsertedResultId()
	if err != nil {
		return -1, err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, cardRes := range result.Result.CardsRes {
		err = u.cardsResultsRepo.InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
		if err != nil {
			return -1, err
		}
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = u.modulesResultsRepo.InsertResultToModule(result.ModuleId, insertedResId)
	if err != nil {
		return -1, err
	}
	return insertedResId, nil
}

func (u *UseCase) InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) ([]int, error) {
	insertedResIds := []int{}
	for _, modulesRes := range result.Modules {
		u.resultsMutex.Lock()
		defer u.resultsMutex.Unlock()

		time, err := time.Parse(time.RFC3339, modulesRes.Result.Time)
		if err != nil {
			return []int{}, err
		}

		err = u.resultRepo.InsertResult(entity.Result{Owner: modulesRes.Result.Owner,
			Type: modulesRes.Result.Type,
			Time: time})
		if err != nil {
			return []int{}, err
		}

		insertedResId, err := u.resultRepo.GetLastInsertedResultId()
		if err != nil {
			return []int{}, err
		}

		u.cardsResultsMutex.Lock()
		defer u.cardsResultsMutex.Unlock()

		for _, cardRes := range modulesRes.Result.CardsRes {
			err = u.cardsResultsRepo.InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
			if err != nil {
				return []int{}, err
			}
		}
		insertedResIds = append(insertedResIds, insertedResId)

		u.categoryModulesResultsMutex.Lock()
		defer u.categoryModulesResultsMutex.Unlock()

		lastInsertedResId, err := u.categoryModulesResultsRepo.GetLastInsertedResId()
		if err != nil {
			return []int{}, err
		}

		err = u.categoryModulesResultsRepo.InsertCategoryModule(lastInsertedResId+1, result.CategoryId, modulesRes.ModuleId, insertedResId)
		if err != nil {
			return []int{}, err
		}
	}
	return insertedResIds, nil
}

func (u *UseCase) DeleteModuleResult(resultId int) error {
	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err := u.cardsResultsRepo.DeleteCardsToResult(resultId)
	if err != nil {
		return err
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = u.modulesResultsRepo.DeleteResultToModule(resultId)
	if err != nil {
		return err
	}

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	return u.resultRepo.DeleteResultById(resultId)
}

func (u *UseCase) DeleteCategoryResultById(categoryResultId int) error {
	categoryRes, err := u.categoryModulesResultsRepo.GetCategoryResById(categoryResultId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = u.categoryModulesResultsRepo.DeleteResultById(categoryResultId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = u.cardsResultsRepo.DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = u.resultRepo.DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) DeleteResultByModuleId(moduleId int) error {
	modulesRes, err := u.modulesResultsRepo.GetResultsToModule(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range modulesRes {
		err := u.cardsResultsRepo.DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = u.modulesResultsRepo.DeleteResultToModule(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = u.resultRepo.DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) DeleteResultByCategoryId(categoryId int) error {
	categoryRes, err := u.categoryModulesResultsRepo.GetCategoryResById(categoryId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = u.categoryModulesResultsRepo.DeleteAllToCategory(categoryId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = u.cardsResultsRepo.DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = u.resultRepo.DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) DeleteModuleResFromCategories(moduleId int) error {
	resultsIds, err := u.categoryModulesResultsRepo.GetResultsByModuleId(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := u.cardsResultsRepo.DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = u.resultRepo.DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return u.categoryModulesResultsRepo.DeleteModulesFromCategories(moduleId)
}

func (u *UseCase) DeleteModuleResFromCategory(categoryId, moduleId int) error {
	resultsIds, err := u.categoryModulesResultsRepo.GetResultsByCategoryAndModule(categoryId, moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := u.cardsResultsRepo.DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = u.resultRepo.DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return u.categoryModulesResultsRepo.DeleteModulesFromCategory(categoryId, moduleId)
}
