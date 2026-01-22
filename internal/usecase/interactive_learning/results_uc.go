package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/uow"
	"interactive_learning/internal/usecase"
	"time"
)

func (u *UseCase) GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error) {
	categoriesRes, err := u.categoryModulesResultsRepoRead.GetCategoriesResByOwner(userId)
	if err != nil {
		return nil, nil, u.errorsMapper.DBErrorToApp(err)
	}

	modulesRes, err := u.modulesResultsRepoRead.GetModulesResByOwner(userId)
	if err != nil {
		return nil, nil, u.errorsMapper.DBErrorToApp(err)
	}

	return categoriesRes, modulesRes, nil
}

func (u *UseCase) GetModuleResultById(resultId int) (entity.ModuleResult, error) {
	moduleResult, err := u.modulesResultsRepoRead.GetModulesResultById(resultId)
	if err != nil {
		return entity.ModuleResult{}, u.errorsMapper.DBErrorToApp(err)
	}
	return moduleResult, nil
}

func (u *UseCase) GetCardsResultById(resultId int) ([]entity.CardsResult, error) {
	cardsResult, err := u.cardsResultsRepoRead.GetCardsResultById(resultId)
	if err != nil {
		return []entity.CardsResult{}, u.errorsMapper.DBErrorToApp(err)
	}
	return cardsResult, nil
}

func (u *UseCase) GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error) {
	modulesResults, err := u.modulesResultsRepoRead.GetResultsToModuleOwner(moduleId, userId)
	if err != nil {
		return []entity.ModuleResult{}, u.errorsMapper.DBErrorToApp(err)
	}
	return modulesResults, nil
}

func (u *UseCase) GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	categoryResults, err := u.categoryModulesResultsRepoRead.GetResultsByCategoryOwner(categoryId, userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, u.errorsMapper.DBErrorToApp(err)
	}
	return categoryResults, nil
}

func (u *UseCase) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	categoryResult, err := u.categoryModulesResultsRepoRead.GetCategoryResById(categoryResultsId)
	if err != nil {
		return entity.CategoryModulesResult{}, u.errorsMapper.DBErrorToApp(err)
	}
	return categoryResult, nil
}

func (u *UseCase) InsertModuleResult(result httputils.InsertModuleResultReq) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	time, err := time.Parse(time.DateTime, result.Time)
	if err != nil {
		return -1, usecase.NewInternalError(err)
	}

	err = uow.GetResultsRepoWriter().InsertResult(entity.Result{
		Type: result.Result.Type})
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}

	insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, cardRes := range result.Result.CardsRes {
		err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
		if err != nil {
			return -1, u.errorsMapper.DBErrorToApp(err)
		}
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().InsertResultToModule(result.ModuleId, insertedResId, result.Owner, time)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return -1, usecase.NewInternalError(err)
	}

	return insertedResId, nil
}

func (u *UseCase) InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) (int, []int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, []int{}, usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	insertedResIds := []int{}

	lastInsertedResId, err := uow.GetCategoryModulesResultsRepoReader().GetLastInsertedResId()
	if err != nil {
		return -1, []int{}, err
	}
	newInsertResultId := lastInsertedResId + 1

	u.resultsMutex.Lock()
	u.cardsResultsMutex.Lock()
	u.categoryModulesResultsMutex.Lock()
	defer func() {
		u.cardsResultsMutex.Unlock()
		u.resultsMutex.Unlock()
		u.categoryModulesResultsMutex.Unlock()
	}()

	for _, modulesRes := range result.Modules {
		err = uow.GetResultsRepoWriter().InsertResult(entity.Result{
			Type: modulesRes.Result.Type})
		if err != nil {
			return -1, []int{}, u.errorsMapper.DBErrorToApp(err)
		}

		insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
		if err != nil {
			return -1, []int{}, u.errorsMapper.DBErrorToApp(err)
		}

		for _, cardRes := range modulesRes.Result.CardsRes {
			err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
			if err != nil {
				return -1, []int{}, u.errorsMapper.DBErrorToApp(err)
			}
		}
		insertedResIds = append(insertedResIds, insertedResId)

		time, err := time.Parse(time.DateTime, result.Time)
		if err != nil {
			return -1, []int{}, usecase.NewInternalError(err)
		}

		err = uow.GetCategoryModulesResultsRepoWriter().InsertCategoryModule(newInsertResultId, result.CategoryId, modulesRes.ModuleId, insertedResId, result.Owner, time)
		if err != nil {
			return -1, []int{}, u.errorsMapper.DBErrorToApp(err)
		}
	}

	if err := uow.Commit(); err != nil {
		return -1, []int{}, usecase.NewInternalError(err)
	}

	return newInsertResultId, insertedResIds, nil
}

func (u *UseCase) DeleteModuleResult(resultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(resultId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) DeleteCategoryResultById(categoryResultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryResultId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteResultById(categoryResultId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) deleteResultByModuleId(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return usecase.NewInternalError(errors.New("uow is null"))
	}

	modulesRes, err := uow.GetModulesResultsRepoReader().GetResultsToModule(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range modulesRes {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}
	return nil
}

func (u *UseCase) deleteResultByCategoryId(categoryId int, uow uow.UnitOfWork) error {
	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.categoryModulesResultsMutex.Lock()
	u.cardsResultsMutex.Lock()
	u.resultsMutex.Lock()
	defer func() {
		u.categoryModulesResultsMutex.Unlock()
		u.cardsResultsMutex.Unlock()
		u.resultsMutex.Unlock()
	}()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteAllToCategory(categoryId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}
	return nil
}

func (u *UseCase) deleteModuleResFromCategories(moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByModuleId(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	u.resultsMutex.Lock()
	u.categoryModulesResultsMutex.Lock()
	defer func() {
		u.cardsResultsMutex.Unlock()
		u.resultsMutex.Unlock()
		u.categoryModulesResultsMutex.Unlock()
	}()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	if err = uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategories(moduleId); err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	return nil
}

func (u *UseCase) deleteModuleResFromCategory(categoryId, moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByCategoryAndModule(categoryId, moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	u.resultsMutex.Lock()
	u.categoryModulesResultsMutex.Lock()
	defer func() {
		u.cardsResultsMutex.Unlock()
		u.resultsMutex.Unlock()
		u.categoryModulesResultsMutex.Unlock()
	}()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	if err = uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategory(categoryId, moduleId); err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	return nil
}
