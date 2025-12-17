package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (cr *CategoryRepo) GetCategoriesToUser(userId int) ([]entity.Category, error) {
	rows, err := cr.db.Query("SELECT * FROM categories WHERE owner_id = $1", userId)
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
	row := cr.db.QueryRow("SELECT * FROM categories WHERE id = $1", id)

	category := entity.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.OwnerId)
	if err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (cr *CategoryRepo) GetLastInsertedCategoryId() (int, error) {
	row := cr.db.QueryRow("SELECT MAX(id) FROM categories")

	var last_id int
	err := row.Scan(&last_id)
	if err != nil {
		return -1, err
	}
	return last_id, nil
}

func (cr *CategoryRepo) InsertCategory(category entity.CategoryToCreate) error {
	result, err := cr.db.Exec("INSERT INTO categories(name, owner_id) "+
		"VALUES($1, $2)", category.Name, category.OwnerId)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert category error")
	}
	return nil
}

func (cr *CategoryRepo) DeleteCategory(id int) error {
	result, err := cr.db.Exec("DELETE FROM categories "+
		"WHERE id = $1", id)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("delete category error")
	}
	return nil
}
