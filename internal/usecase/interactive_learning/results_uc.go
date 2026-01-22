package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/uow"
	"time"
)

func (u *UseCase) GetResultsByOwner(userId int) ([]entity.CategoryModulesResult, []entity.ModuleResult, error) {
	categoriesRes, err := u.categoryModulesResultsRepoRead.GetCategoriesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	modulesRes, err := u.modulesResultsRepoRead.GetModulesResByOwner(userId)
	if err != nil {
		return []entity.CategoryModulesResult{}, []entity.ModuleResult{}, err
	}

	return categoriesRes, modulesRes, nil
}

func (u *UseCase) GetModuleResultById(resultId int) (entity.ModuleResult, error) {
	return u.modulesResultsRepoRead.GetModulesResultById(resultId)
}

func (u *UseCase) GetCardsResultById(resultId int) ([]entity.CardsResult, error) {
	return u.cardsResultsRepoRead.GetCardsResultById(resultId)
}

func (u *UseCase) GetResultsToModuleId(moduleId, userId int) ([]entity.ModuleResult, error) {
	return u.modulesResultsRepoRead.GetResultsToModuleOwner(moduleId, userId)
}

func (u *UseCase) GetResultsByCategoryId(categoryId, userId int) ([]entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepoRead.GetResultsByCategoryOwner(categoryId, userId)
}

func (u *UseCase) GetCategoryResById(categoryResultsId int) (entity.CategoryModulesResult, error) {
	return u.categoryModulesResultsRepoRead.GetCategoryResById(categoryResultsId)
}

func (u *UseCase) InsertModuleResult(result httputils.InsertModuleResultReq) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, err
	}
	defer uow.Rollback()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	time, err := time.Parse(time.DateTime, result.Time)
	if err != nil {
		return -1, err
	}

	err = uow.GetResultsRepoWriter().InsertResult(entity.Result{
		Type: result.Result.Type})
	if err != nil {
		return -1, err
	}

	insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
	if err != nil {
		return -1, err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, cardRes := range result.Result.CardsRes {
		err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
		if err != nil {
			return -1, err
		}
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().InsertResultToModule(result.ModuleId, insertedResId, result.Owner, time)
	if err != nil {
		return -1, err
	}

	if err = uow.Commit(); err != nil {
		return -1, err
	}

	return insertedResId, nil
}

func (u *UseCase) InsertCategoryResult(result httputils.InsertCategoryModulesResultReq) (int, []int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, []int{}, err
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
			return -1, []int{}, err
		}

		insertedResId, err := uow.GetResultsRepoReader().GetLastInsertedResultId()
		if err != nil {
			return -1, []int{}, err
		}

		for _, cardRes := range modulesRes.Result.CardsRes {
			err = uow.GetCardsResultsRepoWriter().InsertCardResult(insertedResId, cardRes.CardId, cardRes.Result)
			if err != nil {
				return -1, []int{}, err
			}
		}
		insertedResIds = append(insertedResIds, insertedResId)

		time, err := time.Parse(time.DateTime, result.Time)
		if err != nil {
			return -1, []int{}, err
		}

		err = uow.GetCategoryModulesResultsRepoWriter().InsertCategoryModule(newInsertResultId, result.CategoryId, modulesRes.ModuleId, insertedResId, result.Owner, time)
		if err != nil {
			return -1, []int{}, err
		}
	}

	if err := uow.Commit(); err != nil {
		return -1, []int{}, err
	}

	return newInsertResultId, insertedResIds, nil
}

func (u *UseCase) DeleteModuleResult(resultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
	if err != nil {
		return err
	}

	u.modulesResultsMutex.Lock()
	defer u.modulesResultsMutex.Unlock()

	err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(resultId)
	if err != nil {
		return err
	}

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
	if err != nil {
		return err
	}

	return uow.Commit()
}

func (u *UseCase) DeleteCategoryResultById(categoryResultId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return err
	}
	defer uow.Rollback()

	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryResultId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteResultById(categoryResultId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return uow.Commit()
}

func (u *UseCase) deleteResultByModuleId(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return errors.New("uow is null")
	}

	modulesRes, err := uow.GetModulesResultsRepoReader().GetResultsToModule(moduleId)
	if err != nil {
		return err
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
			return err
		}

		err = uow.GetModulesResultsRepoWriter().DeleteResultToModule(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) deleteResultByCategoryId(categoryId int, uow uow.UnitOfWork) error {
	categoryRes, err := uow.GetCategoryModulesResultsRepoReader().GetCategoryResById(categoryId)
	if err != nil {
		return err
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	err = uow.GetCategoryModulesResultsRepoWriter().DeleteAllToCategory(categoryId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, moduleRes := range categoryRes.Modules {
		err = uow.GetCardsResultsRepoWriter().DeleteCardsToResult(moduleRes.Result.Id)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(moduleRes.Result.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *UseCase) deleteModuleResFromCategories(moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByModuleId(moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategories(moduleId)
}

func (u *UseCase) deleteModuleResFromCategory(categoryId, moduleId int, uow uow.UnitOfWork) error {
	resultsIds, err := uow.GetCategoryModulesResultsRepoReader().GetResultsByCategoryAndModule(categoryId, moduleId)
	if err != nil {
		return err
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	u.resultsMutex.Lock()
	defer u.resultsMutex.Unlock()

	for _, resultId := range resultsIds {
		err := uow.GetCardsResultsRepoWriter().DeleteCardsToResult(resultId)
		if err != nil {
			return err
		}

		err = uow.GetResultsRepoWriter().DeleteResultById(resultId)
		if err != nil {
			return err
		}
	}

	u.categoryModulesResultsMutex.Lock()
	defer u.categoryModulesResultsMutex.Unlock()

	return uow.GetCategoryModulesResultsRepoWriter().DeleteModulesFromCategory(categoryId, moduleId)
}
