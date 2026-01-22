package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	myerrors "interactive_learning/internal/errors"
)

func (u *UseCase) GetModulesWithSimilarName(name string, limit, offset, userId int) ([]entity.Module, error) {
	modules, err := u.moduleRepoRead.GetModulesWithSimilarName(name, limit, offset)
	if err != nil {
		return nil, err
	}

	publicModules := []entity.Module{}
	for _, module := range modules {
		if module.Type == entity.PrivateModule && module.OwnerId == userId {
			publicModules = append(publicModules, module)
		} else if module.Type == entity.PublicModule {
			publicModules = append(publicModules, module)
		}
	}

	return publicModules, nil
}

func (u *UseCase) GetModulesByUser(ownerId int, withCards bool, userId int) ([]entity.Module, error) {
	modules, err := u.moduleRepoRead.GetModulesByUser(ownerId)
	if err != nil {
		return []entity.Module{}, err
	}

	if ownerId != userId {
		publicModules := []entity.Module{}
		for _, module := range modules {
			if module.Type == entity.PublicModule {
				publicModules = append(publicModules, module)
			}
		}
		modules = publicModules
	}

	if !withCards {
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

func (u *UseCase) GetModuleById(moduleId, userId int) (entity.Module, error) {
	module, err := u.moduleRepoRead.GetModuleById(moduleId)
	if err != nil {
		return entity.Module{}, err
	}

	if module.Type == entity.PrivateModule && userId != module.OwnerId {
		return entity.Module{}, myerrors.NewNotAvailableError("module", moduleId)
	}

	cards, err := u.cardsRepoRead.GetCardsByModule(moduleId)
	if err != nil {
		return entity.Module{}, err
	}
	module.Cards = cards
	return module, nil
}

func (u *UseCase) GetModulesByIds(modulesIds []int, isFull bool, userId int) ([]entity.Module, error) {
	modules := []entity.Module{}

	for _, moduleId := range modulesIds {
		module, err := u.moduleRepoRead.GetModuleById(moduleId)
		if err != nil {
			return []entity.Module{}, err
		}

		if module.Type == entity.PrivateModule && userId != module.OwnerId {
			return []entity.Module{}, myerrors.NewNotAvailableError("module", moduleId)
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

func (u *UseCase) GetPopularModules(limit, offset int) ([]entity.PopularModule, error) {
	return u.moduleRepoRead.GetPopularModules(limit, offset)
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

func (u *UseCase) UpdateModuleType(moduleId, newType, userId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.moduleMutex.Lock()
	u.categoryMutex.Lock()
	u.categoryModulesMutex.Lock()
	defer func() {
		u.moduleMutex.Unlock()
		u.categoryMutex.Unlock()
		u.categoryModulesMutex.Unlock()
	}()

	module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
	if err != nil {
		return err
	}
	if module.OwnerId != userId {
		return errors.New("unaccessable module")
	}

	if module.Type == newType {
		return errors.New("module already has this type")
	}

	err = uow.GetModuleRepoWriter().UpdateModuleType(moduleId, newType)
	if err != nil {
		return err
	}

	switch newType {
	case entity.PrivateModule:
		categories, err := uow.GetCategoryModulesRepoReader().GetCategoriesContainsModule(moduleId)
		if err != nil {
			return err
		}
		for _, category := range categories {
			if category.Type == entity.PublicCategory {
				category.Type = entity.PrivateCategory
			}
			category.Type++
			err = uow.GetCategoryRepoWriter().UpdateCategoryType(category.Id, category.Type)
			if err != nil {
				return err
			}
		}
	case entity.PublicModule:
		categories, err := uow.GetCategoryModulesRepoReader().GetCategoriesContainsModule(moduleId)
		if err != nil {
			return err
		}
		for _, category := range categories {
			err = uow.GetCategoryRepoWriter().TurnDownCategoryType(category.Id)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("Invalid type")
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
	u.selectedMutex.Lock()
	defer func() {
		u.moduleMutex.Unlock()
		u.selectedMutex.Unlock()
	}()

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

	err = uow.GetSelectedRepoWriter().DeleteAllToModule(moduleId)
	if err != nil {
		return err
	}

	err = uow.GetModuleRepoWriter().DeleteModule(moduleId)
	if err != nil {
		return err
	}

	return uow.Commit()
}
