package interactivelearning

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/repo/persistent"
	"interactive_learning/internal/utils/tokengenerator"
	"sync"
)

type UseCase struct {
	userRepo            repo.UsersRepo
	tokenStorage        repo.TokenStorage
	cardRepo            repo.CardRepo
	moduleRepo          repo.ModuleRepo
	categoryRepo        repo.CategoryRepo
	categoryModulesRepo repo.CategoryModulesRepo

	usersMutex           sync.Mutex
	cardMutex            sync.Mutex
	moduleMutex          sync.Mutex
	categoryMutex        sync.Mutex
	categoryModulesMutex sync.Mutex
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
	var ids []int
	var err error
	var curId int

	for _, card := range cards.Cards {
		err = u.cardRepo.InsertCard(entity.Card{Id: card.Id, ParentModule: cards.ParentModule, Term: card.Term, Definition: card.Definition})
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

func (u *UseCase) DeleteCard(cardId int) error {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	return u.cardRepo.DeleteCard(cardId)
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
		cards, err := u.cardRepo.GetCardsByModule(userId)
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

func (u *UseCase) DeleteModule(moduleId int) error {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

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
	for _, module_id := range category.Modules {
		if err = u.InsertModuleToCategory(new_id, module_id); err != nil {
			return -1, err
		}
	}
	return new_id, nil
}

func (u *UseCase) DeleteCategory(id int) error {
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

func (u *UseCase) InsertModuleToCategory(categoryId, moduleId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.InsertModuleToCategory(categoryId, moduleId)
}

func (u *UseCase) DeleteModuleFromCategory(categoryId, moduleId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.DeleteModuleFromCategory(categoryId, moduleId)
}

func (u *UseCase) DeleteAllModulesFromCategory(categoryId int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.DeleteAllModulesFromCategory(categoryId)
}
