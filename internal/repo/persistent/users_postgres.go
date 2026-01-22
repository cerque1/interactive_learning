package persistent

import (
	"database/sql"
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

func (u *UsersRepo) GetUsersWithSimilarName(name string, limit, offset int) ([]entity.User, error) {
	name = "%" + name + "%"
	rows, err := u.psql.Query("SELECT users.id, users.name FROM users WHERE name LIKE $1 LIMIT $2 OFFSET $3", name, limit, offset)
	if err != nil {
		return []entity.User{}, repo.NewDBError("users", "select", err)
	}

	users := []entity.User{}
	for rows.Next() {
		u := entity.User{}
		if err = rows.Scan(&u.Id, &u.Name); err != nil {
			return []entity.User{}, repo.NewDBError("users", "select", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *UsersRepo) GetUserByLogin(login string) (entity.User, error) {
	row := u.psql.QueryRow("select * from users where login = $1", login)

	user := entity.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Name, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, repo.NoSuchRecordToSelect
		}
		return entity.User{}, repo.NewDBError("users", "select", err)
	}

	return user, nil
}

func (u *UsersRepo) GetUserInfoById(userId int) (entity.User, error) {
	row := u.psql.QueryRow("select id, login, name from users where id = $1", userId)

	user := entity.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, repo.NoSuchRecordToSelect
		}
		return entity.User{}, repo.NewDBError("users", "select", err)
	}

	return user, nil
}

func (u *UsersRepo) IsContainsLogin(login string) (bool, error) {
	row := u.psql.QueryRow("select count(*) from users where login = $1", login)

	var count int
	err := row.Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repo.NoSuchRecordToSelect
		}
		return false, repo.NewDBError("users", "select", err)
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
		return repo.NewDBError("users", "insert", err)
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return repo.InsertRecordError
	}
	return nil
}
