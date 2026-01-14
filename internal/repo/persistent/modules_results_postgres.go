package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
	"time"
)

type ModulesResultsRepo struct {
	psql repo.PSQL
}

func NewModulesResultsRepo(psql repo.PSQL) *ModulesResultsRepo {
	return &ModulesResultsRepo{psql: psql}
}

func (mrr *ModulesResultsRepo) GetModulesResultById(resultId int) (entity.ModuleResult, error) {
	row := mrr.psql.QueryRow("SELECT modules_res.module_id, modules_res.result_id, modules_res.time, modules_res.owner, results.type FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.result_id = $1", resultId)

	moduleRes := entity.ModuleResult{}
	if err := row.Scan(&moduleRes.ModuleId,
		&moduleRes.Result.Id,
		&moduleRes.Time,
		&moduleRes.Owner,
		&moduleRes.Result.Type); err != nil {
		return entity.ModuleResult{}, err
	}
	return moduleRes, nil
}

func (mrr *ModulesResultsRepo) GetModulesResByOwner(ownerId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.module_id, modules_res.time, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.\"owner\" = $1", ownerId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{}
		mr.Owner = ownerId
		err := rows.Scan(&mr.ModuleId,
			&mr.Time,
			&mr.Result.Id,
			&mr.Result.Type)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) GetResultsToModuleOwner(moduleId, ownerId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.time, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.module_id = $1 AND modules_res.owner = $2", moduleId, ownerId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{ModuleId: moduleId, Owner: ownerId}
		err := rows.Scan(
			&mr.Time,
			&mr.Result.Id,
			&mr.Result.Type)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) GetResultsToModule(moduleId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.\"owner\", modules_res.time, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.module_id = $1", moduleId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{ModuleId: moduleId}
		err := rows.Scan(&mr.Owner,
			&mr.Time,
			&mr.Result.Id,
			&mr.Result.Type)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) InsertResultToModule(moduleId, resultId, ownerId int, time time.Time) error {
	res, err := mrr.psql.Exec("INSERT INTO modules_res(module_id, result_id, owner, time) "+
		"VALUES($1, $2, $3, $4)", moduleId, resultId, ownerId, time)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("insert module result error")
	}
	return nil
}

func (mrr *ModulesResultsRepo) DeleteResultsToModule(moduleId int) error {
	_, err := mrr.psql.Exec("DELETE FROM modules_res WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (mrr *ModulesResultsRepo) DeleteResultToModule(resultId int) error {
	_, err := mrr.psql.Exec("DELETE FROM modules_res WHERE result_id = $1", resultId)
	if err != nil {
		return err
	}
	return nil
}
