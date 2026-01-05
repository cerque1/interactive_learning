package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type ModulesResultsRepo struct {
	db *sql.DB
}

func NewModulesResultsRepo(db *sql.DB) *ModulesResultsRepo {
	return &ModulesResultsRepo{db: db}
}

func (mrr *ModulesResultsRepo) GetModulesResByOwner(ownerId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.db.Query("SELECT modules_res.module_id, results.id, results.\"owner\", results.\"type\", results.\"time\", results.correct, results.incorrect "+
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
			&mr.Result.Time,
			&mr.Result.Correct,
			&mr.Result.Incorrect)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) GetResultsToModule(moduleId int) ([]entity.ModuleResult, error) {
	rows, err := mrr.db.Query("SELECT modules_res.module_id, results.id, results.\"owner\", results.\"type\", results.\"time\", results.correct, results.incorrect "+
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
			&mr.Result.Time,
			&mr.Result.Correct,
			&mr.Result.Incorrect)
		if err != nil {
			return []entity.ModuleResult{}, err
		}
		modules_results = append(modules_results, mr)
	}

	return modules_results, err
}

func (mrr *ModulesResultsRepo) InsertResultToModule(moduleId, resultId int) error {
	res, err := mrr.db.Exec("INSERT INTO modules_res(module_id, result_id) "+
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
	_, err := mrr.db.Exec("DELETE FROM modules_res WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (mrr *ModulesResultsRepo) DeleteResultToModule(moduleId, resultId int) error {
	_, err := mrr.db.Exec("DELETE FROM modules_res WHERE module_id = $1 AND result_id = $2", moduleId, resultId)
	if err != nil {
		return err
	}
	return nil
}
