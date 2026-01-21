package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type SelectedRepo struct {
	psql repo.PSQL
}

func NewSelectedRepo(psql repo.PSQL) *SelectedRepo {
	return &SelectedRepo{psql: psql}
}

func (sr *SelectedRepo) GetAllSelectedModulesByUser(userId int) ([]entity.Module, error) {
	rows, err := sr.psql.Query("SELECT modules.* FROM selected_modules INNER JOIN modules ON selected_modules.module_id = modules.id "+
		"WHERE user_id = $1", userId)
	if err != nil {
		return []entity.Module{}, err
	}

	modules := []entity.Module{}
	for rows.Next() {
		m := entity.Module{}
		if err := rows.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type); err != nil {
			return []entity.Module{}, err
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (sr *SelectedRepo) GetAllSelectedCategoriesByUser(userId int) ([]entity.Category, error) {
	rows, err := sr.psql.Query("SELECT categories.* FROM selected_categories INNER JOIN categories ON selected_categories.category_id = categories.id "+
		"WHERE user_id = $1", userId)
	if err != nil {
		return []entity.Category{}, err
	}

	categories := []entity.Category{}
	for rows.Next() {
		c := entity.Category{}
		if err := rows.Scan(&c.Id, &c.Name, &c.OwnerId, &c.Type); err != nil {
			return []entity.Category{}, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (sr *SelectedRepo) GetUsersCountToSelectedModule(moduleId int) (int, error) {
	row := sr.psql.QueryRow("SELECT COUNT(DISTINCT user_id) FROM selected_modules")

	var count int
	if err := row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (sr *SelectedRepo) GetUsersCountToSelectedCategory(categoryId int) (int, error) {
	row := sr.psql.QueryRow("SELECT COUNT(DISTINCT user_id) FROM selected_categories")

	var count int
	if err := row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (sr *SelectedRepo) InsertSelectedModuleToUser(userId, moduleId int) error {
	result, err := sr.psql.Exec("INSERT INTO selected_modules(user_id, module_id) VALUES($1, $2)", userId, moduleId)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert selected module error")
	}
	return nil
}

func (sr *SelectedRepo) InsertSelectedCategoryToUser(userId, categoryId int) error {
	result, err := sr.psql.Exec("INSERT INTO selected_categories(user_id, category_id) VALUES($1, $2)", userId, categoryId)
	if err != nil {
		return err
	} else if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert selected category error")
	}
	return nil
}

func (sr *SelectedRepo) DeleteAllToModule(moduleId int) error {
	_, err := sr.psql.Exec("DELETE FROM selected_modules WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SelectedRepo) DeleteAllToCategory(categoryId int) error {
	_, err := sr.psql.Exec("DELETE FROM selected_categories WHERE category_id = $1", categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SelectedRepo) DeleteModuleToUser(userId, moduleId int) error {
	_, err := sr.psql.Exec("DELETE FROM selected_modules WHERE user_id = $1 AND module_id = $2", userId, moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SelectedRepo) DeleteCategoryToUser(userId, categoryId int) error {
	_, err := sr.psql.Exec("DELETE FROM selected_categories WHERE user_id = $1 AND category_id = $2", userId, categoryId)
	if err != nil {
		return err
	}
	return nil
}
