package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type ModulesRepo struct {
	psql repo.PSQL
}

func NewModulesRepo(psql repo.PSQL) *ModulesRepo {
	return &ModulesRepo{psql: psql}
}

func (mr *ModulesRepo) GetModulesWithSimilarName(name string, limit, offset int) ([]entity.Module, error) {
	name = "%" + name + "%"
	rows, err := mr.psql.Query("SELECT modules.id, modules.name, modules.owner_id, modules.type FROM modules WHERE name LIKE $1 LIMIT $2 OFFSET $3", name, limit, offset)
	if err != nil {
		return []entity.Module{}, repo.NewDBError("modules", "select", err)
	}

	modules := []entity.Module{}
	for rows.Next() {
		m := entity.Module{}
		if err = rows.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type); err != nil {
			return []entity.Module{}, repo.NewDBError("modules", "select", err)
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (mr *ModulesRepo) GetModulesByUser(userId int) ([]entity.Module, error) {
	rows, err := mr.psql.Query("SELECT * FROM modules WHERE owner_id = $1", userId)
	if err != nil {
		return []entity.Module{}, repo.NewDBError("modules", "select", err)
	}

	modules := []entity.Module{}
	for rows.Next() {
		m := entity.Module{}
		err = rows.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type)
		if err != nil {
			return []entity.Module{}, repo.NewDBError("modules", "select", err)
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (cr *ModulesRepo) GetModuleById(moduleId int) (entity.Module, error) {
	row := cr.psql.QueryRow("SELECT * FROM modules WHERE id = $1", moduleId)
	m := entity.Module{}
	err := row.Scan(&m.Id, &m.Name, &m.OwnerId, &m.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Module{}, repo.NoSuchRecordToSelect
		}
		return entity.Module{}, repo.NewDBError("modules", "select", err)
	}
	return m, nil
}

func (mr *ModulesRepo) GetLastInsertedModuleId() (int, error) {
	row := mr.psql.QueryRow("SELECT MAX(id) FROM modules")
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, repo.NoSuchRecordToSelect
		}
		return -1, repo.NewDBError("modules", "select", err)
	}
	return id, nil
}

func (mr *ModulesRepo) GetModuleOwnerId(moduleId int) (int, error) {
	row := mr.psql.QueryRow("SELECT owner_id FROM modules WHERE id = $1", moduleId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, repo.NoSuchRecordToSelect
		}
		return -1, repo.NewDBError("modules", "select", err)
	}
	return id, nil
}

func (mr *ModulesRepo) GetPopularModules(limit, offset int) ([]entity.PopularModule, error) {
	rows, err := mr.psql.Query("SELECT modules.*, COUNT(DISTINCT modules_res.owner) as count "+
		"FROM modules INNER JOIN modules_res ON modules.id = modules_res.module_id "+
		"WHERE modules.type = 0 AND time >= NOW() - INTERVAL '7 days' "+
		"GROUP BY modules.id, modules_res.owner "+
		"ORDER BY count "+
		"LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return []entity.PopularModule{}, repo.NewDBError("modules", "select", err)
	}

	modules := []entity.PopularModule{}
	for rows.Next() {
		m := entity.PopularModule{}
		err = rows.Scan(&m.Mod.Id, &m.Mod.Name, &m.Mod.OwnerId, &m.Mod.Type, &m.Count)
		if err != nil {
			return []entity.PopularModule{}, repo.NewDBError("modules", "select", err)
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (mr *ModulesRepo) InsertModule(module entity.ModuleToCreate) error {
	result, err := mr.psql.Exec("INSERT INTO modules(name, owner_id, type) "+
		"VALUES($1, $2, $3)", module.Name, module.OwnerId, module.Type)
	if err != nil {
		return repo.NewDBError("modules", "insert", err)
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return repo.InsertRecordError
	}
	return nil
}

func (mr *ModulesRepo) RenameModule(moduleId int, newName string) error {
	result, err := mr.psql.Exec("UPDATE modules "+
		"SET name = $1 "+
		"WHERE id = $2", newName, moduleId)
	if err != nil {
		return repo.NewDBError("modules", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (mr *ModulesRepo) UpdateModuleType(moduleId, newType int) error {
	result, err := mr.psql.Exec("UPDATE modules "+
		"SET type = $1 "+
		"WHERE id = $2", newType, moduleId)
	if err != nil {
		return repo.NewDBError("modules", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (mr *ModulesRepo) DeleteModule(moduleId int) error {
	result, err := mr.psql.Exec("DELETE FROM modules WHERE id = $1", moduleId)
	if err != nil {
		return repo.NewDBError("modules", "delete", err)
	} else if count, _ := result.RowsAffected(); count < 0 {
		return repo.NoSuchRecordToDelete
	}
	return nil
}
