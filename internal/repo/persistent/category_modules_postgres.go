package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type CategoryModulesRepo struct {
	psql repo.PSQL
}

func NewCategoryModulesRepo(psql repo.PSQL) *CategoryModulesRepo {
	return &CategoryModulesRepo{psql: psql}
}

func (cmr *CategoryModulesRepo) GetModulesToCategory(categoryId int) ([]entity.Module, error) {
	rows, err := cmr.psql.Query("SELECT id, name, owner_id, type FROM category_modules LEFT JOIN modules ON modules.id = category_modules.module_id "+
		"WHERE category_id = $1", categoryId)
	if err != nil {
		return []entity.Module{}, err
	}

	modules := []entity.Module{}
	for rows.Next() {
		m := entity.Module{}
		err := rows.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type)
		if err != nil {
			return []entity.Module{}, err
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (cmr *CategoryModulesRepo) InsertModulesToCategory(categoryId, moduleId int) error {
	result, err := cmr.psql.Exec("INSERT INTO category_modules(category_id, module_id) "+
		"VALUES($1, $2)", categoryId, moduleId)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert module to category error")
	}
	return nil
}

func (cmr *CategoryModulesRepo) DeleteModuleFromCategory(categoryId, moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_modules "+
		"WHERE category_id = $1 AND module_id = $2", categoryId, moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesRepo) DeleteAllModulesFromCategory(categoryId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_modules "+
		"WHERE category_id = $1", categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesRepo) DeleteModuleFromCategories(moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_modules "+
		"WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}
