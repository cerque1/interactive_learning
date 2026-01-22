package persistent

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
	"time"
)

type CategoryModulesResultsRepo struct {
	psql repo.PSQL
}

func NewCategoryModulesResultsRepo(psql repo.PSQL) *CategoryModulesResultsRepo {
	return &CategoryModulesResultsRepo{psql: psql}
}

func (cmr *CategoryModulesResultsRepo) GetCategoriesResByOwner(ownerId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.psql.Query("SELECT category_res.category_result_id, category_res.category_id, category_res.module_id, category_res.time, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_res.\"owner\" = $1 "+
		"ORDER BY category_res.category_result_id", ownerId)
	if err != nil {
		return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
	}

	categoryRes := map[int]entity.CategoryModulesResult{}
	first := true
	var tempResId int
	var tempCategoryId int
	var tempTime string
	oneCategoryRes := entity.CategoryModulesResult{Owner: ownerId}
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempResId,
			&tempCategoryId,
			&moduleRes.ModuleId,
			&tempTime,
			&moduleRes.Result.Id,
			&moduleRes.Result.Type)
		if err != nil {
			return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
		}

		if first {
			oneCategoryRes.CategoryResultId = tempResId
			oneCategoryRes.CategoryId = tempCategoryId
			parseTime, err := time.Parse(time.RFC3339, tempTime)
			if err != nil {
				return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
			}
			oneCategoryRes.Time = parseTime

			first = false
		}

		if oneCategoryRes.CategoryResultId != tempResId {
			categoryRes[oneCategoryRes.CategoryResultId] = oneCategoryRes

			parseTime, err := time.Parse(time.RFC3339, tempTime)
			if err != nil {
				return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
			}
			oneCategoryRes = entity.CategoryModulesResult{CategoryResultId: tempResId, CategoryId: tempCategoryId, Owner: ownerId, Time: parseTime}
		}

		oneCategoryRes.Modules = append(oneCategoryRes.Modules, moduleRes)
	}
	if !first {
		categoryRes[oneCategoryRes.CategoryResultId] = oneCategoryRes
	}

	categoryResArray := []entity.CategoryModulesResult{}
	for _, value := range categoryRes {
		categoryResArray = append(categoryResArray, value)
	}

	return categoryResArray, nil
}

func (cmr *CategoryModulesResultsRepo) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	rows, err := cmr.psql.Query("SELECT category_res.category_id, category_res.\"owner\", category_res.module_id, category_res.time, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_result_id = $1", categoryResultsId)
	if err != nil {
		return entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
	}

	categoryRes := entity.CategoryModulesResult{CategoryResultId: categoryResultsId}
	first := true
	var tempCategoryId int
	var tempOwner int
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempCategoryId,
			&tempOwner,
			&moduleRes.ModuleId,
			&categoryRes.Time,
			&moduleRes.Result.Id,
			&moduleRes.Result.Type)
		if err != nil {
			return entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
		}
		if first {
			categoryRes.CategoryId = tempCategoryId
			categoryRes.Owner = tempOwner
			first = false
		}

		categoryRes.Modules = append(categoryRes.Modules, moduleRes)
	}

	return categoryRes, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryOwner(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.psql.Query("SELECT category_res.category_result_id, category_res.module_id, category_res.time, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_id = $1 AND category_res.owner = $2 "+
		"ORDER BY category_res.category_result_id", categoryId, userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
	}

	categoryRes := map[int]entity.CategoryModulesResult{}
	first := true
	var tempResId int
	var tempTime string
	oneCategoryRes := entity.CategoryModulesResult{CategoryId: categoryId, Owner: userId}
	for rows.Next() {
		moduleRes := entity.ModuleResult{}
		err = rows.Scan(&tempResId,
			&moduleRes.ModuleId,
			&tempTime,
			&moduleRes.Result.Id,
			&moduleRes.Result.Type)
		if err != nil {
			return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
		}

		if first {
			oneCategoryRes.CategoryResultId = tempResId
			parseTime, err := time.Parse(time.RFC3339, tempTime)
			if err != nil {
				return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
			}
			oneCategoryRes.Time = parseTime

			first = false
		}

		if oneCategoryRes.CategoryResultId != tempResId {
			categoryRes[oneCategoryRes.CategoryResultId] = oneCategoryRes

			parseTime, err := time.Parse(time.RFC3339, tempTime)
			if err != nil {
				return []entity.CategoryModulesResult{}, repo.NewDBError("category_res", "select", err)
			}
			oneCategoryRes = entity.CategoryModulesResult{CategoryResultId: tempResId, CategoryId: categoryId, Owner: userId, Time: parseTime}
		}

		oneCategoryRes.Modules = append(oneCategoryRes.Modules, moduleRes)
	}
	if !first {
		categoryRes[oneCategoryRes.CategoryResultId] = oneCategoryRes
	}

	categoryResArray := []entity.CategoryModulesResult{}
	for _, value := range categoryRes {
		categoryResArray = append(categoryResArray, value)
	}

	return categoryResArray, nil
}

func (cmr *CategoryModulesResultsRepo) GetLastInsertedResId() (int, error) {
	row := cmr.psql.QueryRow("SELECT COALESCE(MAX(category_result_id), 0) AS max_id FROM category_res")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, repo.NewDBError("category_res", "select", err)
	}
	return id, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByModuleId(moduleId int) ([]int, error) {
	rows, err := cmr.psql.Query("SELECT result_id FROM category_res WHERE module_id = $1", moduleId)
	if err != nil {
		return []int{}, repo.NewDBError("category_res", "select", err)
	}

	ids := []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return []int{}, repo.NewDBError("category_res", "select", err)
		}

		ids = append(ids, id)
	}
	return ids, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryAndModule(categoryId, moduleId int) ([]int, error) {
	rows, err := cmr.psql.Query("SELECT result_id FROM category_res WHERE category_id = $1 AND module_id = $2", categoryId, moduleId)
	if err != nil {
		return []int{}, repo.NewDBError("category_res", "select", err)
	}

	ids := []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return []int{}, repo.NewDBError("category_res", "select", err)
		}

		ids = append(ids, id)
	}
	return ids, nil
}

func (cmr *CategoryModulesResultsRepo) InsertCategoryModule(categoryResultId, categoryId, moduleId, resultId, ownerId int, time time.Time) error {
	res, err := cmr.psql.Exec("INSERT INTO category_res(category_result_id, category_id, module_id, result_id, owner, time) "+
		"VALUES($1, $2, $3, $4, $5, $6)", categoryResultId, categoryId, moduleId, resultId, ownerId, time)
	if err != nil {
		return repo.InsertRecordError
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return repo.InsertRecordError
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteModulesFromCategories(moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE module_id = $1", moduleId)
	if err != nil {
		return repo.NoSuchRecordToDelete
	}
	return nil
}

// rename method
func (cmr *CategoryModulesResultsRepo) DeleteModulesFromCategory(categoryId, moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_id = $1 AND module_id = $2", categoryId, moduleId)
	if err != nil {
		return repo.NoSuchRecordToDelete
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteAllToCategory(categoryId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_id = $1", categoryId)
	if err != nil {
		return repo.NoSuchRecordToDelete
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteResultById(categoryResultId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_result_id = $1", categoryResultId)
	if err != nil {
		return repo.NoSuchRecordToDelete
	}
	return nil
}
