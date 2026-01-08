package persistent

import (
	"errors"
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
	rows, err := cmr.psql.Query("SELECT category_res.category_result_id, category_res.category_id, category_res.module_id, results.* "+
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
			&timeStr)
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
	rows, err := cmr.psql.Query("SELECT category_res.category_id, category_res.module_id, results.* "+
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
			&timeStr)
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

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryOwner(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.psql.Query("SELECT category_res.category_result_id, category_res.module_id, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_id = $1 AND results.owner = $2 "+
		"ORDER BY category_res.category_result_id", categoryId, userId)
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
			&timeStr)
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

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryId(categoryId int) ([]entity.CategoryModulesResult, error) {
	rows, err := cmr.psql.Query("SELECT category_res.category_result_id, category_res.module_id, results.* "+
		"FROM category_res INNER JOIN results ON category_res.result_id = results.id "+
		"WHERE category_id = $1 "+
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
			&timeStr)
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
	row := cmr.psql.QueryRow("SELECT MAX(category_result_id) FROM category_res")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByModuleId(moduleId int) ([]int, error) {
	rows, err := cmr.psql.Query("SELECT result_id FROM category_res WHERE module_id = $1", moduleId)
	if err != nil {
		return []int{}, err
	}

	ids := []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return []int{}, err
		}

		ids = append(ids, id)
	}
	return ids, nil
}

func (cmr *CategoryModulesResultsRepo) GetResultsByCategoryAndModule(categoryId, moduleId int) ([]int, error) {
	rows, err := cmr.psql.Query("SELECT result_id FROM category_res WHERE category_id = $1 AND module_id = $2", categoryId, moduleId)
	if err != nil {
		return []int{}, err
	}

	ids := []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return []int{}, err
		}

		ids = append(ids, id)
	}
	return ids, nil
}

func (cmr *CategoryModulesResultsRepo) InsertCategoryModule(categoryResultId, categoryId, moduleId, result_id int) error {
	res, err := cmr.psql.Exec("INSERT INTO category_res(category_result_id, category_id, module_id, result_id) "+
		"VALUES($1, $2)", categoryResultId, categoryId, moduleId, result_id)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("insert category result error")
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteModulesFromCategories(moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteModulesFromCategory(categoryId, moduleId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_id = $1 AND module_id = $2", categoryId, moduleId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteAllToCategory(categoryId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_id = $1", categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CategoryModulesResultsRepo) DeleteResultById(categoryResultId int) error {
	_, err := cmr.psql.Exec("DELETE FROM category_res WHERE category_result_id = $1", categoryResultId)
	if err != nil {
		return err
	}
	return nil
}
