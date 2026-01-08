package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type CategoryRepo struct {
	psql repo.PSQL
}

func NewCategoryRepo(psql repo.PSQL) *CategoryRepo {
	return &CategoryRepo{psql: psql}
}

func (cr *CategoryRepo) GetCategoriesToUser(userId int) ([]entity.Category, error) {
	rows, err := cr.psql.Query("SELECT * FROM categories WHERE owner_id = $1", userId)
	if err != nil {
		return []entity.Category{}, nil
	}

	categories := []entity.Category{}
	for rows.Next() {
		c := entity.Category{}
		err = rows.Scan(&c.Id, &c.Name, &c.OwnerId)
		if err != nil {
			return []entity.Category{}, nil
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (cr *CategoryRepo) GetCategoryById(id int) (entity.Category, error) {
	row := cr.psql.QueryRow("SELECT * FROM categories WHERE id = $1", id)

	category := entity.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.OwnerId)
	if err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (cr *CategoryRepo) GetLastInsertedCategoryId() (int, error) {
	row := cr.psql.QueryRow("SELECT MAX(id) FROM categories")

	var last_id int
	err := row.Scan(&last_id)
	if err != nil {
		return -1, err
	}
	return last_id, nil
}

func (cr *CategoryRepo) GetCategoryOwnerId(categoryId int) (int, error) {
	row := cr.psql.QueryRow("SELECT owner_id FROM categories WHERE id = $1", categoryId)

	var ownerId int
	err := row.Scan(&ownerId)
	if err != nil {
		return -1, err
	}
	return ownerId, nil
}

func (cr *CategoryRepo) InsertCategory(category entity.CategoryToCreate) error {
	result, err := cr.psql.Exec("INSERT INTO categories(name, owner_id) "+
		"VALUES($1, $2)", category.Name, category.OwnerId)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert category error")
	}
	return nil
}

func (cr *CategoryRepo) RenameCategory(categoryId int, newName string) error {
	result, err := cr.psql.Exec("UPDATE categories "+
		"SET name = $1 "+
		"WHERE id = $2", newName, categoryId)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("rename category error")
	}
	return nil
}

func (cr *CategoryRepo) DeleteCategory(id int) error {
	_, err := cr.psql.Exec("DELETE FROM categories "+
		"WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
