package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type ModulesResultsRepo struct {
	psql repo.PSQL
}

func NewModulesResultsRepo(psql repo.PSQL) *ModulesResultsRepo {
	return &ModulesResultsRepo{psql: psql}
}

func (mrr *ModulesResultsRepo) GetModulesResByOwner(ownerId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.module_id, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE results.\"owner\" = $1", ownerId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{}
		err := rows.Scan(&mr.ModuleId,
			&mr.Result.Id,
			&mr.Result.Owner,
			&mr.Result.Type,
			&mr.Result.Time)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) GetResultsToModuleOwner(moduleId, ownerId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.module_id, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.module_id = $1 AND results.owner = $2", moduleId, ownerId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{}
		err := rows.Scan(&mr.ModuleId,
			&mr.Result.Id,
			&mr.Result.Owner,
			&mr.Result.Type,
			&mr.Result.Time)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) GetResultsToModule(moduleId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.psql.Query("SELECT modules_res.module_id, results.* "+
		"FROM modules_res INNER JOIN results ON modules_res.result_id = results.id "+
		"WHERE modules_res.module_id = $1", moduleId)
	if err != nil {
		return []entity.ModuleResult{}, err
	}

	modules_results := []entity.ModuleResult{}
	for rows.Next() {
		mr := entity.ModuleResult{}
		err := rows.Scan(&mr.ModuleId,
			&mr.Result.Id,
			&mr.Result.Owner,
			&mr.Result.Type,
			&mr.Result.Time)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) InsertResultToModule(moduleId, resultId int) error {
	res, err := mrr.psql.Exec("INSERT INTO modules_res(module_id, result_id) "+
		"VALUES($1, $2)", moduleId, resultId)
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
	_, err := mrr.psql.Exec("DELETE FROM modules_res WHERE result_id = $2", resultId)
	if err != nil {
		return err
	}
	return nil
}
