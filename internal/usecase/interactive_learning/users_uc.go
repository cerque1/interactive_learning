package interactivelearning

import "interactive_learning/internal/entity"

func (u *UseCase) GetUserByLogin(login string) (entity.User, error) {
	return u.usersRepoRead.GetUserByLogin(login)
}

func (u *UseCase) GetUsersWithSimilarName(name string, limit, offset int) ([]entity.User, error) {
	return u.usersRepoRead.GetUsersWithSimilarName(name, limit, offset)
}

func (u *UseCase) GetUserInfoById(ownerId int, isFull bool, userId int) (entity.User, error) {
	user, err := u.usersRepoRead.GetUserInfoById(ownerId)
	if err != nil {
		return entity.User{}, err
	}

	if !isFull {
		return user, nil
	}

	modules, err := u.GetModulesByUser(ownerId, true, userId)
	if err != nil {
		return entity.User{}, err
	}
	user.Modules = modules

	categories, err := u.GetCategoriesToUser(ownerId, true, userId)
	if err != nil {
		return entity.User{}, err
	}
	user.Categories = categories

	return user, nil
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
