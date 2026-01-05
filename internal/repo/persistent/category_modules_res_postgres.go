package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
	"time"
)

type CategoryModulesResultsRepo struct {
	db *sql.DB
}

func NewCategoryModulesResultsRepo(db *sql.DB) *CategoryModulesResultsRepo {
	return &CategoryModulesResultsRepo{db: db}
}

func (cmr *CategoryModulesResultsRepo) GetCategoriesResByOwner(ownerId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.db.Query("SELECT category_res.category_result_id, category_res.category_id, category_res.module_id, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE results.\"owner\" = $1 "+
		"ORDER BY category_res.category_result_id", ownerId)
	if err != nil {
		return []entity.CategoryModulesResult{}, err
	}

	categoryRes := map[int]entity.CategoryModulesResult{}
	first := true
	var tempResId int
	var timeStr string
	oneCategoryRes := entity.CategoryModulesResult{}
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempResId,
			&oneCategoryRes.CategoryId,
			&moduleRes.ModuleId,
			&moduleRes.Result.Id,
			&moduleRes.Result.Owner,
			&moduleRes.Result.Type,
			&timeStr,
			&moduleRes.Result.Correct,
			&moduleRes.Result.Incorrect)
		if err != nil {
			return []entity.CategoryModulesResult{}, err
		}

		if first {
			oneCategoryRes.CategoryResultId = tempResId
			first = false
		}

		if oneCategoryRes.CategoryId != tempResId {
			categoryRes[oneCategoryRes.CategoryId] = oneCategoryRes
			oneCategoryRes = entity.CategoryModulesResult{CategoryResultId: tempResId}
		}

		parseTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return []entity.CategoryModulesResult{}, err
		}
		moduleRes.Result.Time = parseTime

		oneCategoryRes.Modules = append(oneCategoryRes.Modules, moduleRes)
	}
	categoryRes[oneCategoryRes.CategoryId] = oneCategoryRes

	categoryResArray := []entity.CategoryModulesResult{}
	for _, value := range categoryRes {
		categoryResArray = append(categoryResArray, value)
	}

	return categoryResArray, nil
}

func (cmr *CategoryModulesResultsRepo) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	rows, err := cmr.db.Query("SELECT category_res.category_id, category_res.module_id, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_result_id = $1", categoryResultsId)
	if err != nil {
		return entity.CategoryModulesResult{}, err
	}

	categoryRes := entity.CategoryModulesResult{CategoryResultId: categoryResultsId}
	first := true
	var tempCategoryId int
	var timeStr string
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempCategoryId,
			&moduleRes.ModuleId,
			&moduleRes.Result.Id,
			&moduleRes.Result.Owner,
			&moduleRes.Result.Type,
			&timeStr,
			&moduleRes.Result.Correct,
			&moduleRes.Result.Incorrect)
		if err != nil {
			return entity.CategoryModulesResult{}, err
		}
		if first {
			categoryRes.CategoryId = tempCategoryId
			first = false
		}

		parseTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return entity.CategoryModulesResult{}, err
		}
		moduleRes.Result.Time = parseTime

		categoryRes.Modules = append(categoryRes.Modules, moduleRes)
	}

	return categoryRes, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryId(categoryId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.db.Query("SELECT category_res.category_result_id, category_res.module_id, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_id = $1"+
		"ORDER BY category_res.category_result_id", categoryId)
	if err != nil {
		return []entity.CategoryModulesResult{}, err
	}

	categoryRes := map[int]entity.CategoryModulesResult{}
	first := true
	var tempResId int
	var timeStr string
	oneCategoryRes := entity.CategoryModulesResult{CategoryId: categoryId}
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempResId,
			&moduleRes.ModuleId,
			&moduleRes.Result.Id,
			&moduleRes.Result.Owner,
			&moduleRes.Result.Type,
			&timeStr,
			&moduleRes.Result.Correct,
			&moduleRes.Result.Incorrect)
		if err != nil {
			return []entity.CategoryModulesResult{}, err
		}

		if first {
			oneCategoryRes.CategoryResultId = tempResId
			first = false
		}

		if oneCategoryRes.CategoryId != tempResId {
			categoryRes[oneCategoryRes.CategoryId] = oneCategoryRes
			oneCategoryRes = entity.CategoryModulesResult{CategoryId: categoryId, CategoryResultId: tempResId}
		}

		parseTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return []entity.CategoryModulesResult{}, err
		}
		moduleRes.Result.Time = parseTime

		oneCategoryRes.Modules = append(oneCategoryRes.Modules, moduleRes)
	}
	categoryRes[oneCategoryRes.CategoryId] = oneCategoryRes

	categoryResArray := []entity.CategoryModulesResult{}
	for _, value := range categoryRes {
		categoryResArray = append(categoryResArray, value)
	}

	return categoryResArray, nil
}

func (cmr *CategoryModulesResultsRepo) GetLastInsertedResId() (int, error) {
	row := cmr.db.QueryRow("SELECT MAX(category_result_id) FROM category_res")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (cmr *CategoryModulesResultsRepo) InsertCategoryModule(categoryResultId, categoryId, moduleId int) error {
	res, err := cmr.db.Exec("INSERT INTO category_res(category_result_id, category_id, module_id) "+
		"VALUES($1, $2)", categoryResultId, categoryId, moduleId)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("insert category result error")
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteModulesFromCategory(moduleId int) error {
	_, err := cmr.db.Exec("DELETE FROM category_res WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteAllToCategory(categoryId int) error {
	_, err := cmr.db.Exec("DELETE FROM category_res WHERE category_id = $1", categoryId)
	if err != nil {
		return err
	}
	return nil
}
