package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/uow"
	"interactive_learning/internal/usecase"
)

func (u *UseCase) GetCategoriesWithSimilarName(name string, limit, offset, userId int) ([]entity.Category, error) {
	categories, err := u.categoryRepoRead.GetCategoriesWithSimilarName(name, limit, offset)
	if err != nil {
		return nil, u.errorsMapper.DBErrorToApp(err)
	}

	publicCategories := []entity.Category{}
	for _, category := range categories {
		if category.Type == entity.PublicCategory {
			publicCategories = append(publicCategories, category)
		} else if category.Type >= entity.PrivateCategory && category.OwnerId == userId {
			publicCategories = append(publicCategories, category)
		}
	}

	return publicCategories, nil
}

func (u *UseCase) GetCategoriesToUser(ownerId int, isFull bool, userId int) ([]entity.Category, error) {
	categories, err := u.categoryRepoRead.GetCategoriesToUser(ownerId)
	if err != nil {
		return []entity.Category{}, u.errorsMapper.DBErrorToApp(err)
	}

	if userId != ownerId {
		publicCategories := []entity.Category{}
		for _, category := range categories {
			if category.Type == entity.PublicCategory {
				publicCategories = append(publicCategories, category)
			}
		}
		categories = publicCategories
	}

	if !isFull {
		return categories, nil
	}

	for i := range categories {
		modules, err := u.GetModulesToCategory(categories[i].Id, true, userId)
		if err != nil {
			return []entity.Category{}, err
		}

		categories[i].Modules = modules
	}

	return categories, nil
}

func (u *UseCase) GetCategoryById(id int, userId int) (entity.Category, error) {
	category, err := u.categoryRepoRead.GetCategoryById(id)
	if err != nil {
		return entity.Category{}, u.errorsMapper.DBErrorToApp(err)
	}

	if category.Type >= entity.PrivateCategory && category.OwnerId != userId {
		return entity.Category{}, usecase.NewNotAvailableError("category", category.Id)
	}

	modules, err := u.GetModulesToCategory(id, true, userId)
	if err != nil {
		return entity.Category{}, err
	}
	category.Modules = modules

	return category, nil
}

func (u *UseCase) GetPopularCategories(limit, offset int) ([]entity.PopularCategory, error) {
	popularCategories, err := u.categoryRepoRead.GetPopularCategories(limit, offset)
	if err != nil {
		return nil, u.errorsMapper.DBErrorToApp(err)
	}
	return popularCategories, nil
}

func (u *UseCase) InsertCategory(category entity.CategoryToCreate) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.categoryMutex.Lock()
	defer u.categoryMutex.Unlock()

	err := uow.GetCategoryRepoWriter().InsertCategory(category)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}
	new_id, err := uow.GetCategoryRepoReader().GetLastInsertedCategoryId()
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}
	if err = u.insertModulesToCategory(category.OwnerId, new_id, category.Modules, uow); err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, usecase.NewInternalError(err)
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
		return false, u.errorsMapper.DBErrorToApp(err)
	}
	if ownerId != userId {
		return false, nil
	}
	return true, nil
}

func (u *UseCase) RenameCategory(userId, categoryId int, newName string) error {
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

	err = uow.GetCategoryRepoWriter().RenameCategory(categoryId, newName)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) UpdateCategoryType(categoryId, newType, userId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	category, err := uow.GetCategoryRepoReader().GetCategoryById(categoryId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	isOwner, err := u.isCategoryOwner(userId, categoryId, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return usecase.NewNotAvailableError("category", categoryId)
	}

	if category.Type == newType {
		return usecase.NewChangeTypeError("category", errors.New("category already has this type"))
	} else if newType == entity.PublicCategory && category.Type > entity.PrivateCategory {
		return usecase.NewChangeTypeError("category", errors.New("It is impossible to set the type, the category contains a closed module."))
	} else if newType > entity.PrivateCategory || newType < entity.PublicCategory {
		return usecase.NewChangeTypeError("category", errors.New("Invalid type"))
	}

	err = uow.GetCategoryRepoWriter().UpdateCategoryType(categoryId, newType)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) DeleteCategory(userId int, id int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.categoryMutex.Lock()
	u.selectedMutex.Lock()
	defer func() {
		u.categoryMutex.Unlock()
		u.selectedMutex.Unlock()
	}()

	isOwner, err := u.isCategoryOwner(userId, id, uow)
	if err != nil {
		return err
	} else if !isOwner {
		return usecase.NewNotAvailableError("category", id)
	}

	err = u.deleteAllModulesFromCategory(id, uow)
	if err != nil {
		return err
	}

	err = u.deleteResultByCategoryId(id, uow)
	if err != nil {
		return err
	}

	err = uow.GetSelectedRepoWriter().DeleteAllToCategory(id)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	err = uow.GetCategoryRepoWriter().DeleteCategory(id)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}
