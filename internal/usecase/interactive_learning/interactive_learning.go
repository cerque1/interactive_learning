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

func (u *UseCase) GetUserInfoById(user_id int, is_full bool) (entity.User, error) {
	user, err := u.userRepo.GetUserInfoById(user_id)
	if err != nil {
		return entity.User{}, err
	}

	if !is_full {
		return user, nil
	}

	modules, err := u.GetModulesWithCardsByUser(user_id)
	if err != nil {
		return entity.User{}, err
	}
	user.Modules = modules

	categories, err := u.GetCategoriesToUser(user_id, true)
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
	new_user, err := u.userRepo.GetUserByLogin(user.Login)
	if err != nil {
		return -1, err
	}
	return new_user.Id, nil
}

func (u *UseCase) GetCardById(card_id int) (entity.Card, error) {
	return u.cardRepo.GetCardById(card_id)
}

func (u *UseCase) GetCardsByModule(module_id int) ([]entity.Card, error) {
	return u.cardRepo.GetCardsByModule(module_id)
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

func (u *UseCase) InsertCards(cards []entity.Card) ([]int, error) {
	var ids []int
	var err error
	var cur_id int

	for _, card := range cards {
		err = u.cardRepo.InsertCard(card)
		if err != nil {
			return []int{}, err
		}
		cur_id, err = u.cardRepo.GetLastInsertedCardId()
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, cur_id)
	}
	return ids, nil
}

func (u *UseCase) DeleteCard(card_id int) error {
	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	return u.cardRepo.DeleteCard(card_id)
}

func (u *UseCase) GetModulesByUser(user_id int) ([]entity.Module, error) {
	return u.moduleRepo.GetModulesByUser(user_id)
}

func (u *UseCase) GetModulesWithCardsByUser(user_id int) ([]entity.Module, error) {
	modules, err := u.moduleRepo.GetModulesByUser(user_id)
	if err != nil {
		return []entity.Module{}, err
	}

	for i := range modules {
		cards, err := u.cardRepo.GetCardsByModule(user_id)
		if err != nil {
			return []entity.Module{}, err
		}
		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) GetModuleById(module_id int) (entity.Module, error) {
	module, err := u.moduleRepo.GetModuleById(module_id)
	if err != nil {
		return entity.Module{}, err
	}
	cards, err := u.cardRepo.GetCardsByModule(module_id)
	if err != nil {
		return entity.Module{}, err
	}
	module.Cards = cards
	return module, nil
}

func (u *UseCase) InsertModule(module entity.Module) (int, []int, error) {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	err := u.moduleRepo.InsertModule(module)
	if err != nil {
		return -1, []int{}, err
	}
	insert_ids, err := u.InsertCards(module.Cards)
	if err != nil {
		return -1, []int{}, err
	}
	id, err := u.moduleRepo.GetLastInsertedModuleId()
	if err != nil {
		return -1, []int{}, err
	}
	return id, insert_ids, nil
}

func (u *UseCase) DeleteModule(module_id int) error {
	u.moduleMutex.Lock()
	defer u.moduleMutex.Unlock()

	return u.moduleRepo.DeleteModule(module_id)
}

func (u *UseCase) GetCategoriesToUser(user_id int, is_full bool) ([]entity.Category, error) {
	categories, err := u.categoryRepo.GetCategoriesToUser(user_id)
	if err != nil {
		return []entity.Category{}, err
	}

	if !is_full {
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

func (u *UseCase) InsertCategory(category entity.Category) (int, error) {
	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	err := u.categoryRepo.InsertCategory(category)
	if err != nil {
		return -1, err
	}
	return u.categoryRepo.GetLastInsertedCategoryId()
}

func (u *UseCase) DeleteCategory(id int) error {
	return u.categoryRepo.DeleteCategory(id)
}

func (u *UseCase) GetModulesToCategory(category_id int, is_full bool) ([]entity.Module, error) {
	modules, err := u.categoryModulesRepo.GetModulesToCategory(category_id)
	if err != nil {
		return []entity.Module{}, err
	}

	if !is_full {
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

func (u *UseCase) InsertModuleToCategory(category_id, module_id int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.InsertModuleToCategory(category_id, module_id)
}

func (u *UseCase) DeleteModuleFromCategory(category_id, module_id int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.DeleteModuleFromCategory(category_id, module_id)
}

func (u *UseCase) DeleteAllModulesFromCategory(category_id int) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	return u.categoryModulesRepo.DeleteAllModulesFromCategory(category_id)
}
