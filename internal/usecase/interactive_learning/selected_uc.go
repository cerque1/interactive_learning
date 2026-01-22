package interactivelearning

import "interactive_learning/internal/entity"

func (u *UseCase) GetAllSelectedModulesByUser(userId int) ([]entity.Module, error) {
	modules, err := u.selectedRepoRead.GetAllSelectedModulesByUser(userId)
	if err != nil {
		return []entity.Module{}, err
	}

	publicModules := []entity.Module{}
	for _, module := range modules {
		if module.OwnerId == userId || module.Type == entity.PublicModule {
			publicModules = append(publicModules, module)
		}
	}
	return publicModules, nil
}

func (u *UseCase) GetAllSelectedCategoriesByUser(userId int) ([]entity.Category, error) {
	categories, err := u.selectedRepoRead.GetAllSelectedCategoriesByUser(userId)
	if err != nil {
		return []entity.Category{}, err
	}

	publicCategories := []entity.Category{}
	for _, category := range categories {
		if category.OwnerId == userId || category.Type == entity.PublicCategory {
			publicCategories = append(publicCategories, category)
		}
	}
	return publicCategories, nil
}

func (u *UseCase) GetUsersCountToSelectedModule(moduleId int) (int, error) {
	usersCount, err := u.selectedRepoRead.GetUsersCountToSelectedModule(moduleId)
	if err != nil {
		return -1, err
	}
	return usersCount, nil
}

func (u *UseCase) GetUsersCountToSelectedCategory(categoryId int) (int, error) {
	usersCount, err := u.selectedRepoRead.GetUsersCountToSelectedCategory(categoryId)
	if err != nil {
		return -1, err
	}
	return usersCount, nil
}

func (u *UseCase) InsertSelectedModuleToUser(userId, moduleId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.GetSelectedRepoWriter().InsertSelectedModuleToUser(userId, moduleId); err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) InsertSelectedCategoryToUser(userId, categoryId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.GetSelectedRepoWriter().InsertSelectedCategoryToUser(userId, categoryId); err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteModuleToUser(userId, moduleId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.GetSelectedRepoWriter().DeleteModuleToUser(userId, moduleId); err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteCategoryToUser(userId, categoryId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.GetSelectedRepoWriter().DeleteCategoryToUser(userId, categoryId); err != nil {
		return err
	}

	return uow.Commit()
}
