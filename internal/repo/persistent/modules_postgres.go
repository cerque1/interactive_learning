package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type ModulesRepo struct {
	db *sql.DB
}

func NewModulesRepo(db *sql.DB) *ModulesRepo {
	return &ModulesRepo{db: db}
}

func (mr *ModulesRepo) GetModulesByUser(userId int) ([]entity.Module, error) {
	rows, err := mr.db.Query("SELECT * FROM modules WHERE owner_id = $1", userId)
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

func (cr *ModulesRepo) GetModuleById(moduleId int) (entity.Module, error) {
	row := cr.db.QueryRow("SELECT * FROM modules WHERE id = $1", moduleId)
	m := entity.Module{}
	err := row.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type)
	if err != nil {
		return entity.Module{}, err
	}
	return m, nil
}

func (mr *ModulesRepo) GetLastInsertedModuleId() (int, error) {
	row := mr.db.QueryRow("SELECT MAX(id) FROM modules")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (mr *ModulesRepo) GetModuleOwnerId(moduleId int) (int, error) {
	row := mr.db.QueryRow("SELECT owner_id FROM modules WHERE id = $1", moduleId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (mr *ModulesRepo) InsertModule(module entity.ModuleToCreate) error {
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

func (mr *ModulesRepo) DeleteModule(moduleId int) error {
	_, err := mr.db.Exec("DELETE FROM modules WHERE id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}
