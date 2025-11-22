package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type ModulesRepoImpl struct {
	db *sql.DB
}

func (mr *ModulesRepoImpl) GetModulesByUser(user_id int) ([]entity.Module, error) {
	rows, err := mr.db.Query("SELECT * FROM modules WHERE owner_id = $1", user_id)
	if err != nil {
		return []entity.Module{}, err
	}

	modules := []entity.Module{}
	for rows.Next() {
		m := entity.Module{}
		err = rows.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type)
		if err != nil {
			return []entity.Module{}, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (cr *CardsRepoImpl) GetLastInsertedModuleId() (int, error) {
	row := cr.db.QueryRow("SELECT MAX(id) FROM modules")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.ErrUnsupported
	}
	return id, nil
}

func (mr *ModulesRepoImpl) InsertModule(module entity.Module) error {
	result, err := mr.db.Exec("INSERT INTO modules(name, owner_id, type) "+
		"VALUES($1, $2, $3)", module.Name, module.OwnerId, module.Type)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert module error")
	}
	return nil
}

func (mr *ModulesRepoImpl) DeleteModule(module_id int) error {
	result, err := mr.db.Exec("DELETE FROM modules WHERE id = $1", module_id)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert module error")
	}
	return nil
}
