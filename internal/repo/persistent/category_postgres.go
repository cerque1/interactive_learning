package persistent

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type CategoryRepo struct {
	psql repo.PSQL
}

func NewCategoryRepo(psql repo.PSQL) *CategoryRepo {
	return &CategoryRepo{psql: psql}
}

func (cr *CategoryRepo) GetCategoriesWithSimilarName(name string, limit, offset int) ([]entity.Category, error) {
	name = "%" + name + "%"
	rows, err := cr.psql.Query("SELECT categories.id, categories.name, categories.owner_id, categories.type FROM categories WHERE name LIKE $1 LIMIT $2 OFFSET $3", name, limit, offset)
	if err != nil {
		return []entity.Category{}, repo.NewDBError("categories", "select", err)
	}

	categories := []entity.Category{}
	for rows.Next() {
		c := entity.Category{}
		if err = rows.Scan(&c.Id, &c.Name, &c.OwnerId, &c.Type); err != nil {
			return []entity.Category{}, repo.NewDBError("categories", "select", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (cr *CategoryRepo) GetCategoriesToUser(userId int) ([]entity.Category, error) {
	rows, err := cr.psql.Query("SELECT * FROM categories WHERE owner_id = $1", userId)
	if err != nil {
		return []entity.Category{}, repo.NewDBError("categories", "select", err)
	}

	categories := []entity.Category{}
	for rows.Next() {
		c := entity.Category{}
		err = rows.Scan(&c.Id, &c.Name, &c.OwnerId, &c.Type)
		if err != nil {
			return []entity.Category{}, repo.NewDBError("categories", "select", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (cr *CategoryRepo) GetCategoryById(id int) (entity.Category, error) {
	row := cr.psql.QueryRow("SELECT * FROM categories WHERE id = $1", id)
	category := entity.Category{}
	err := row.Scan(&category.Id, &category.Name, &category.OwnerId, &category.Type)
	if err != nil {
		return entity.Category{}, repo.NoSuchRecordToSelect
	}
	return category, nil
}

func (cr *CategoryRepo) GetLastInsertedCategoryId() (int, error) {
	row := cr.psql.QueryRow("SELECT MAX(id) FROM categories")

	var last_id int
	err := row.Scan(&last_id)
	if err != nil {
		return -1, repo.NewDBError("categories", "select", err)
	}
	return last_id, nil
}

func (cr *CategoryRepo) GetCategoryOwnerId(categoryId int) (int, error) {
	row := cr.psql.QueryRow("SELECT owner_id FROM categories WHERE id = $1", categoryId)

	var ownerId int
	err := row.Scan(&ownerId)
	if err != nil {
		return -1, repo.NewDBError("categories", "select", err)
	}
	return ownerId, nil
}

func (cr *CategoryRepo) GetPopularCategories(limit, offset int) ([]entity.PopularCategory, error) {
	rows, err := cr.psql.Query("SELECT categories.*, COUNT(DISTINCT category_res.owner) as count "+
		"FROM categories INNER JOIN category_res ON categories.id = category_res.category_id "+
		"WHERE categories.type = 0 AND time >= NOW() - INTERVAL '7 days' "+
		"GROUP BY categories.id, category_res.owner "+
		"ORDER BY count "+
		"LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return []entity.PopularCategory{}, repo.NewDBError("categories", "select", err)
	}

	categories := []entity.PopularCategory{}
	for rows.Next() {
		c := entity.PopularCategory{}
		err = rows.Scan(&c.Cat.Id, &c.Cat.Name, &c.Cat.OwnerId, &c.Cat.Type, &c.Count)
		if err != nil {
			return []entity.PopularCategory{}, repo.NewDBError("categories", "select", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (cr *CategoryRepo) InsertCategory(category entity.CategoryToCreate) error {
	result, err := cr.psql.Exec("INSERT INTO categories(name, owner_id, type) "+
		"VALUES($1, $2, $3)", category.Name, category.OwnerId, category.Type)
	if err != nil {
		return repo.NewDBError("categories", "insert", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.InsertRecordError
	}
	return nil
}

func (cr *CategoryRepo) RenameCategory(categoryId int, newName string) error {
	result, err := cr.psql.Exec("UPDATE categories "+
		"SET name = $1 "+
		"WHERE id = $2", newName, categoryId)
	if err != nil {
		return repo.NewDBError("categories", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (cr *CategoryRepo) UpdateCategoryType(categoryId, categoryType int) error {
	result, err := cr.psql.Exec("UPDATE categories "+
		"SET type = $1 "+
		"WHERE id = $2", categoryType, categoryId)
	if err != nil {
		return repo.NewDBError("categories", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (cr *CategoryRepo) TurnDownCategoryType(categoryId int) error {
	result, err := cr.psql.Exec("UPDATE categories "+
		"SET type = type - 1 "+
		"WHERE id = $1", categoryId)
	if err != nil {
		return repo.NewDBError("categories", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (cr *CategoryRepo) DeleteCategory(id int) error {
	result, err := cr.psql.Exec("DELETE FROM categories "+
		"WHERE id = $1", id)
	if err != nil {
		return repo.NewDBError("categories", "delete", err)
	} else if count, _ := result.RowsAffected(); count < 1 {
		return repo.NoSuchRecordToDelete
	}
	return nil
}
