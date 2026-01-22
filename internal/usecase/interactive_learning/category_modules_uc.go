package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/uow"
	"interactive_learning/internal/usecase"
	"log"
	"slices"
)

func (u *UseCase) GetModulesToCategory(categoryId int, isFull bool, userId int) ([]entity.Module, error) {
	category, err := u.categoryRepoRead.GetCategoryById(categoryId)
	if err != nil {
		return []entity.Module{}, u.errorsMapper.DBErrorToApp(err)
	}

	if category.Type >= entity.PrivateCategory && category.OwnerId != userId {
		return []entity.Module{}, usecase.NewNotAvailableError("category", categoryId)
	}

	modules, err := u.categoryModulesRepoRead.GetModulesToCategory(categoryId)
	if err != nil {
		return []entity.Module{}, u.errorsMapper.DBErrorToApp(err)
	}

	if !isFull {
		return modules, nil
	}

	for i := range modules {
		cards, err := u.cardsRepoRead.GetCardsByModule(modules[i].Id)
		if err != nil {
			return []entity.Module{}, u.errorsMapper.DBErrorToApp(err)
		}

		modules[i].Cards = cards
	}

	return modules, nil
}

func (u *UseCase) InsertModulesToCategory(userId, categoryId int, modulesIds []int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	err := u.insertModulesToCategory(userId, categoryId, modulesIds, uow)
	if err != nil {
		return err
	}
	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) insertModulesToCategory(userId, categoryId int, modulesIds []int, uow uow.UnitOfWork) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	category, err := uow.GetCategoryRepoReader().GetCategoryById(categoryId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	if category.OwnerId != userId {
		return usecase.NewNotAvailableError("category", categoryId)
	}

	for _, moduleId := range modulesIds {
		if idx := slices.IndexFunc(category.Modules, func(elt entity.Module) bool { return elt.Id == moduleId }); idx >= 0 {
			return usecase.NewAlreadyExistsError("module", moduleId)
		}
	}

	newCategoryType := category.Type
	for _, moduleId := range modulesIds {
		module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		if module.Type == entity.PrivateModule {
			if newCategoryType == entity.PublicCategory {
				newCategoryType = entity.PrivateCategory
			}
			newCategoryType++
			log.Print(newCategoryType)
		}

		err = uow.GetCategoryModulesRepoWriter().InsertModulesToCategory(categoryId, moduleId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	if newCategoryType != category.Type {
		err = uow.GetCategoryRepoWriter().UpdateCategoryType(categoryId, newCategoryType)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}
	return nil
}

func (u *UseCase) DeleteModuleFromCategory(userId, categoryId, moduleId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	isOwner, err := u.isCategoryOwner(userId, categoryId, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return usecase.NewNotAvailableError("category", categoryId)
	}

	err = u.deleteModuleResFromCategory(categoryId, moduleId, uow)
	if err != nil {
		return err
	}

	module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	if module.Type == entity.PrivateModule {
		err = uow.GetCategoryRepoWriter().TurnDownCategoryType(categoryId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	err = uow.GetCategoryModulesRepoWriter().DeleteModuleFromCategory(categoryId, moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) deleteAllModulesFromCategory(categoryId int, uow uow.UnitOfWork) error {
	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	err := uow.GetCategoryModulesRepoWriter().DeleteAllModulesFromCategory(categoryId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	return nil
}

func (u *UseCase) deleteModuleFromCategories(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return usecase.NewInternalError(errors.New("uow is null"))
	}

	u.categoryModulesMutex.Lock()
	defer u.categoryModulesMutex.Unlock()

	module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	if module.Type == entity.PrivateModule {
		categories, err := uow.GetCategoryModulesRepoReader().GetCategoriesContainsModule(moduleId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
		for _, category := range categories {
			err = uow.GetCategoryRepoWriter().TurnDownCategoryType(category.Id)
			if err != nil {
				return u.errorsMapper.DBErrorToApp(err)
			}
		}
	}

	if err = uow.GetCategoryModulesRepoWriter().DeleteModuleFromCategories(moduleId); err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	return nil
}
