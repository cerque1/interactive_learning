package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type UsersRepo struct {
	psql repo.PSQL
}

func NewUsersRepo(psql repo.PSQL) *UsersRepo {
	return &UsersRepo{psql}
}

func (u *UsersRepo) GetUserByLogin(login string) (entity.User, error) {
	row := u.psql.QueryRow("select * from users where login = $1", login)

	user := entity.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Name, &user.PasswordHash)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *UsersRepo) GetUserInfoById(userId int) (entity.User, error) {
	row := u.psql.QueryRow("select id, login, name from users where id = $1", userId)

	user := entity.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Name)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *UsersRepo) IsContainsLogin(login string) (bool, error) {
	row := u.psql.QueryRow("select count(*) from users where login = $1", login)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (u *UsersRepo) InsertUser(user entity.User) error {
	result, err := u.psql.Exec("insert into users(login, name, password_hash) "+
		"values($1, $2, $3)", user.Login, user.Name, user.PasswordHash)

	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert user error")
	}
	return nil
}
