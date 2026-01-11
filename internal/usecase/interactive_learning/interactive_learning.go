package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/repo/persistent"
	"interactive_learning/internal/uow"
	"interactive_learning/internal/utils/tokengenerator"
	"slices"
	"sync"
	"time"
)

type UseCase struct {
	unitOfWorkFactory func() uow.UnitOfWork

	tokenStorage                   repo.TokenStorage
	usersRepoRead                  repo.UsersRepoRead
	cardsRepoRead                  repo.CardRepoRead
	moduleRepoRead                 repo.ModuleRepoRead
	categoryRepoRead               repo.CategoryRepoRead
	categoryModulesRepoRead        repo.CategoryModulesRepoRead
	resultsRepoRead                repo.ResultsRepoRead
	cardsResultsRepoRead           repo.CardsResultsRepoRead
	modulesResultsRepoRead         repo.ModulesResultsRepoRead
	categoryModulesResultsRepoRead repo.CategoryModulesResultsRepoRead

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

func New(unitOfWorkFactory func() uow.UnitOfWork,
	usersRepoRead repo.UsersRepoRead,
	cardsRepoRead repo.CardRepoRead,
	moduleRepoRead repo.ModuleRepoRead,
	categoryRepoRead repo.CategoryRepoRead,
	categoryModulesRepoRead repo.CategoryModulesRepoRead,
	resultsRepoRead repo.ResultsRepoRead,
	cardsResultsRepoRead repo.CardsResultsRepoRead,
	modulesResultsRepoRead repo.ModulesResultsRepoRead,
	categoryModulesResultsRepoRead repo.CategoryModulesResultsRepoRead) *UseCase {

	return &UseCase{unitOfWorkFactory: unitOfWorkFactory,
		tokenStorage:                   persistent.NewTokenStorage(),
		usersRepoRead:                  usersRepoRead,
		cardsRepoRead:                  cardsRepoRead,
		moduleRepoRead:                 moduleRepoRead,
		categoryRepoRead:               categoryRepoRead,
		categoryModulesRepoRead:        categoryModulesRepoRead,
		resultsRepoRead:                resultsRepoRead,
		cardsResultsRepoRead:           cardsResultsRepoRead,
		modulesResultsRepoRead:         modulesResultsRepoRead,
		categoryModulesResultsRepoRead: categoryModulesResultsRepoRead,
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
	return u.usersRepoRead.GetUserByLogin(login)
}

func (u *UseCase) GetUserInfoById(userId int, isFull bool) (entity.User, error) {
	user, err := u.usersRepoRead.GetUserInfoById(userId)
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

func (u *UseCase) getCardOwnerId(cardId int, uow uow.UnitOfWork) (int, error) {
	cardsRepoRead := u.cardsRepoRead
	moduleRepoRead := u.moduleRepoRead
	if uow != nil {
		cardsRepoRead = uow.GetCardRepoReader()
		moduleRepoRead = uow.GetModuleRepoReader()
	}

	parentModule, err := cardsRepoRead.GetParentModuleId(cardId)
	if err != nil {
		return -1, err
	}
	moduleOwner, err := moduleRepoRead.GetModuleOwnerId(parentModule)
	if err != nil {
		return -1, err
	}
	return moduleOwner, nil
}

func (u *UseCase) IsContainsLogin(login string) (bool, error) {
	u.usersMutex.Lock()
	defer u.usersMutex.Unlock()

	return u.usersRepoRead.IsContainsLogin(login)
}

func (u *UseCase) InsertUser(user entity.User) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, err
	}
	defer uow.Rollback()

	u.usersMutex.Lock()
	defer u.usersMutex.Unlock()

	err := uow.GetUsersRepoWriter().InsertUser(user)
	if err != nil {
		return -1, err
	}
	newUser, err := uow.GetUsersRepoReader().GetUserByLogin(user.Login)
	if err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, err
	}

	return newUser.Id, nil
}

func (u *UseCase) GetCardById(cardId int) (entity.Card, error) {
	return u.cardsRepoRead.GetCardById(cardId)
}

func (u *UseCase) GetCardsByModule(moduleId int) ([]entity.Card, error) {
	return u.cardsRepoRead.GetCardsByModule(moduleId)
}

func (u *UseCase) InsertCard(card entity.Card) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, err
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	err := uow.GetCardRepoWriter().InsertCard(card)
	if err != nil {
		return -1, err
	}

	insertedId, err := uow.GetCardRepoReader().GetLastInsertedCardId()
	if err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, err
	}

	return insertedId, nil
}

func (u *UseCase) InsertCards(cards entity.CardsToAdd) ([]int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return []int{}, err
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	var ids []int
	var err error
	var curId int

	for _, card := range cards.Cards {
		err = uow.GetCardRepoWriter().InsertCard(entity.Card{ParentModule: cards.ParentModule, Term: card.Term, Definition: card.Definition})
		if err != nil {
			return []int{}, err
		}
		curId, err = uow.GetCardRepoReader().GetLastInsertedCardId()
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, curId)
	}

	if err = uow.Commit(); err != nil {
		return []int{}, err
	}

	return ids, nil
}

func (u *UseCase) UpdateCard(userId int, card entity.Card) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.getCardOwnerId(card.Id, uow)
	if err != nil {
		return errors.New("bad card id")
	} else if ownerId != userId {
		return errors.New("unaccessable card")
	}

	err = uow.GetCardRepoWriter().UpdateCard(card)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteCard(userId int, cardId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.getCardOwnerId(cardId, uow)
	if err != nil {
		return errors.New("bad card id")
	}

	if ownerId != userId {
		return errors.New("unaccessable card")
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err = uow.GetCardsResultsRepoWriter().DeleteResultsToCard(cardId)
	if err != nil {
		return err
	}

	err = uow.GetCardRepoWriter().DeleteCard(cardId)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) deleteCardsToParentModule(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return errors.New("uow is null")
	}

	module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, card := range module.Cards {
		err = uow.GetCardsResultsRepoWriter().DeleteResultsToCard(card.Id)
		if err != nil {
			return err
		}
	}

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	err = uow.GetCardRepoWriter().DeleteCardsToParentModule(moduleId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) GetModulesByUser(userId int) ([]entity.Module, error) {
	return u.moduleRepoRead.GetModulesByUser(userId)
}

func (u *UseCase) GetModulesWithCardsByUser(userId int) ([]entity.Module, error) {
	modules, err := u.moduleRepoRead.GetModulesByUser(userId)
	if err != nil {
		return []entity.Module{}, err
	}

	for i := range modules {
		cards, err := u.cardsRepoRead.GetCardsByModule(modules[i].Id)
		if err != nil {
			return []entity.Module{}, err
		}
		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) GetModuleById(moduleId int) (entity.Module, error) {
	module, err := u.moduleRepoRead.GetModuleById(moduleId)
	if err != nil {
		return entity.Module{}, err
	}
	cards, err := u.cardsRepoRead.GetCardsByModule(moduleId)
	if err != nil {
		return entity.Module{}, err
	}
	module.Cards = cards
	return module, nil
}

func (u *UseCase) GetModulesByIds(modulesIds []int, isFull bool) ([]entity.Module, error) {
	modules := []entity.Module{}

	for _, moduleId := range modulesIds {
		module, err := u.moduleRepoRead.GetModuleById(moduleId)
		if err != nil {
			return []entity.Module{}, err
		}

		if isFull {
			cards, err := u.cardsRepoRead.GetCardsByModule(moduleId)
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
	return u.moduleRepoRead.GetModuleOwnerId(moduleId)
}

func (u *UseCase) InsertModule(module entity.ModuleToCreate) (int, []int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, []int{}, err
	}
	defer uow.Rollback()

	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	err := uow.GetModuleRepoWriter().InsertModule(module)
	if err != nil {
		return -1, []int{}, err
	}
	insertIds, err := u.InsertCards(entity.CardsToAdd{Cards: module.Cards, ParentModule: module.Id})
	if err != nil {
		return -1, []int{}, err
	}
	id, err := uow.GetModuleRepoReader().GetLastInsertedModuleId()
	if err != nil {
		return -1, []int{}, err
	}

	if err = uow.Commit(); err != nil {
		return -1, []int{}, err
	}

	return id, insertIds, nil
}

func (u *UseCase) RenameModule(userId, moduleId int, newName string) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	ownerId, err := uow.GetModuleRepoReader().GetModuleOwnerId(moduleId)
	if err != nil {
		return errors.New("bad module id")
	}
	if ownerId != userId {
		return errors.New("unaccessable module")
	}

	err = uow.GetModuleRepoWriter().RenameModule(moduleId, newName)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteModule(userId int, moduleId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	ownerId, err := uow.GetModuleRepoReader().GetModuleOwnerId(moduleId)
	if err != nil {
		return errors.New("bad module id")
	}
	if ownerId != userId {
		return errors.New("unaccessable module")
	}

	err = u.deleteCardsToParentModule(moduleId, uow)
	if err != nil {
		return err
	}

	err = u.deleteModuleFromCategories(moduleId, uow)
	if err != nil {
		return err
	}

	err = u.deleteResultByModuleId(moduleId, uow)
	if err != nil {
		return err
	}

	err = u.deleteModuleResFromCategories(moduleId, uow)
	if err != nil {
		return err
	}

	err = uow.GetModuleRepoWriter().DeleteModule(moduleId)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) GetCategoriesToUser(userId int, isFull bool) ([]entity.Category, error) {
	categories, err := u.categoryRepoRead.GetCategoriesToUser(userId)
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
	category, err := u.categoryRepoRead.GetCategoryById(id)
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
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, err
	}
	defer uow.Rollback()

	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	err := uow.GetCategoryRepoWriter().InsertCategory(category)
	if err != nil {
		return -1, err
	}
	new_id, err := uow.GetCategoryRepoReader().GetLastInsertedCategoryId()
	if err != nil {
		return -1, err
	}
	if err = u.insertModulesToCategory(category.OwnerId, new_id, category.Modules, uow); err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, err
	}

	return new_id, nil
}

func (u *UseCase) IsCategoryOwner(userId, categoryId int) (bool, error) {
	return u.isCategoryOwner(userId, categoryId, nil)
}

func (u *UseCase) isCategoryOwner(userId, categoryId int, uow uow.UnitOfWork) (bool, error) {
	categoryRepoRead := u.categoryRepoRead
	if uow != nil {
		categoryRepoRead = uow.GetCategoryRepoReader()
	}

	ownerId, err := categoryRepoRead.GetCategoryOwnerId(categoryId)
	if err != nil {
		return false, errors.New("bad category id")
	}
	if ownerId != userId {
		return false, nil
	}
	return true, nil
}

func (u *UseCase) RenameCategory(userId, categoryId int, newName string) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.isCategoryOwner(userId, categoryId, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	err = uow.GetCategoryRepoWriter().RenameCategory(categoryId, newName)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteCategory(userId int, id int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	isOwner, err := u.isCategoryOwner(userId, id, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	err = u.deleteAllModulesFromCategory(id, uow)
	if err != nil {
		return err
	}

	err = u.deleteResultByCategoryId(id, uow)
	if err != nil {
		return err
	}

	err = uow.GetCategoryRepoWriter().DeleteCategory(id)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) GetModulesToCategory(categoryId int, isFull bool) ([]entity.Module, error) {
	modules, err := u.categoryModulesRepoRead.GetModulesToCategory(categoryId)
	if err != nil {
		return []entity.Module{}, err
	}

	if !isFull {
		return modules, nil
	}

	for i := range modules {
		cards, err := u.cardsRepoRead.GetCardsByModule(modules[i].Id)
		if err != nil {
			return []entity.Module{}, err
		}

		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) InsertModulesToCategory(userId, categoryId int, modulesIds []int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	err := u.insertModulesToCategory(userId, categoryId, modulesIds, uow)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) insertModulesToCategory(userId, categoryId int, modulesIds []int, uow uow.UnitOfWork) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	category, err := uow.GetCategoryRepoReader().GetCategoryById(categoryId)
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
		err := uow.GetCategoryModulesRepoWriter().InsertModulesToCategory(categoryId, moduleId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) DeleteModuleFromCategory(userId, categoryId, moduleId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.isCategoryOwner(userId, categoryId, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return errors.New("unavailable category")
	}

	err = u.deleteModuleResFromCategory(categoryId, moduleId, uow)
	if err != nil {
		return err
	}

	err = uow.GetCategoryModulesRepoWriter().DeleteModuleFromCategory(categoryId, moduleId)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) deleteAllModulesFromCategory(categoryId int, uow uow.UnitOfWork) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return uow.GetCategoryModulesRepoWriter().DeleteAllModulesFromCategory(categoryId)
}

func (u *UseCase) deleteModuleFromCategories(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return errors.New("uow is null")
	}

	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return uow.GetCategoryModulesRepoWriter().DeleteModuleFromCategories(moduleId)
}

func (u *UseCase) GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error) {
	categoriesRes, err := u.categoryModulesResultsRepoRead.GetCategoriesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	modulesRes, err := u.modulesResultsRepoRead.GetModulesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	return categoriesRes, modulesRes, nil
}

func (u *UseCase) GetCardsResultById(resultId int) ([]entity.CardsResult, error) {
	return u.cardsResultsRepoRead.GetCardsResultById(resultId)
}

func (u *UseCase) GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error) {
	return u.modulesResultsRepoRead.GetResultsToModuleOwner(moduleId, userId)
}

func (u *UseCase) GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepoRead.GetResultsByCategoryOwner(categoryId, userId)
}

func (u *UseCase) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepoRead.GetCategoryResById(categoryResultsId)
}

func (u *UseCase) InsertModuleResult(result httputils.InsertModuleResultReq) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, err
	}
	defer uow.Rollback()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	time, err := time.Parse(time.DateTime, result.Result.Time)
	if err != nil {
		return -1, err
	}

	err = uow.GetResultsRepoWriter().InsertResult(entity.Result{Owner: result.Result.Owner,
		Type: result.Result.Type,
		Time: time})
	if err != nil {
		return -1, err
	}

	insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
	if err != nil {
		return -1, err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, cardRes := range result.Result.CardsRes {
		err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
		if err != nil {
			return -1, err
		}
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().InsertResultToModule(result.ModuleId, insertedResId)
	if err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, err
	}

	return insertedResId, nil
}

func (u *UseCase) InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) (int, []int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, []int{}, err
	}
	defer uow.Rollback()

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	insertedResIds := []int{}

	lastInsertedResId, err := uow.GetCategoryModulesResultsRepoReader().GetLastInsertedResId()
	if err != nil {
		return -1, []int{}, err
	}
	newInsertResultId := lastInsertedResId + 1

	for _, modulesRes := range result.Modules {
		u.resultsMutex.Lock()
		defer u.resultsMutex.Unlock()

		time, err := time.Parse(time.DateTime, modulesRes.Result.Time)
		if err != nil {
			return -1, []int{}, err
		}

		err = uow.GetResultsRepoWriter().InsertResult(entity.Result{Owner: modulesRes.Result.Owner,
			Type: modulesRes.Result.Type,
			Time: time})
		if err != nil {
			return -1, []int{}, err
		}

		insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
		if err != nil {
			return -1, []int{}, err
		}

		u.cardsResultsMutex.Lock()
		defer u.cardsResultsMutex.Unlock()

		for _, cardRes := range modulesRes.Result.CardsRes {
			err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
			if err != nil {
				return -1, []int{}, err
			}
		}
		insertedResIds = append(insertedResIds, insertedResId)

		err = uow.GetCategoryModulesResultsRepoWriter().InsertCategoryModule(newInsertResultId, result.CategoryId, modulesRes.ModuleId, insertedResId)
		if err != nil {
			return -1, []int{}, err
		}
	}

	if err := uow.Commit(); err != nil {
		return -1, []int{}, err
	}

	return newInsertResultId, insertedResIds, nil
}

func (u *UseCase) DeleteModuleResult(resultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
	if err != nil {
		return err
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(resultId)
	if err != nil {
		return err
	}

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteCategoryResultById(categoryResultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryResultId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteResultById(categoryResultId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return uow.Commit()
}

func (u *UseCase) deleteResultByModuleId(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return errors.New("uow is null")
	}

	modulesRes, err := uow.GetModulesResultsRepoReader().GetResultsToModule(moduleId)
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
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) deleteResultByCategoryId(categoryId int, uow uow.UnitOfWork) error {
	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteAllToCategory(categoryId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) deleteModuleResFromCategories(moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByModuleId(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategories(moduleId)
}

func (u *UseCase) deleteModuleResFromCategory(categoryId, moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByCategoryAndModule(categoryId, moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategory(categoryId, moduleId)
}
