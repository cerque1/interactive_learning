package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type ResultsRepo struct {
	psql repo.PSQL
}

func NewResultsRepo(psql repo.PSQL) *ResultsRepo {
	return &ResultsRepo{psql: psql}
}

func (rr *ResultsRepo) GetResultsByOwner(ownerId int) ([]entity.Result, error) {
	rows, err := rr.psql.Query("SELECT * FROM results WHERE owner = $1", ownerId)
	if err != nil {
		return []entity.Result{}, err
	}
	defer rows.Close()

	results := []entity.Result{}
	for rows.Next() {
		r := entity.Result{}
		err = rows.Scan(&r.Id,
			&r.Owner,
			&r.Type,
			&r.Time)
		if err != nil {
			return []entity.Result{}, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (rr *ResultsRepo) GetResultById(id int) (entity.Result, error) {
	row := rr.psql.QueryRow("SELECT * FROM results id = $1", id)
	r := entity.Result{}
	err := row.Scan(&r.Id,
		&r.Owner,
		&r.Type,
		&r.Time)
	if err != nil {
		return entity.Result{}, err
	}
	return r, nil
}

func (rr *ResultsRepo) GetLastInsertedResultId() (int, error) {
	row := rr.psql.QueryRow("SELECT MAX(id) FROM results")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (rr *ResultsRepo) InsertResult(result entity.Result) error {
	res, err := rr.psql.Exec("INSERT INTO results(owner, type, time) "+
		"VALUES($1, $2, $3, $4, $5)", result.Owner, result.Type, result.Time)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("insert result error")
	}
	return nil
}

func (rr *ResultsRepo) DeleteResultById(id int) error {
	_, err := rr.psql.Exec("DELETE FROM results WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
