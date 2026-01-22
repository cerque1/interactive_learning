package selected

import (
	errors_mapper "interactive_learning/internal/mappers/errors"
	"interactive_learning/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SelectedRouter struct {
	selectedUC usecase.Selected

	errorsMapper *errors_mapper.ApplicationErrorsMapper
}

func NewSelectedRouter(selectedUC usecase.Selected, errorsMapper *errors_mapper.ApplicationErrorsMapper) *SelectedRouter {
	return &SelectedRouter{selectedUC: selectedUC, errorsMapper: errorsMapper}
}

func (sr *SelectedRouter) GetAllSelectedModulesByUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	modules, err := sr.selectedUC.GetAllSelectedModulesByUser(userId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"selected_modules": modules,
	})
}

func (sr *SelectedRouter) GetAllSelectedCategoriesByUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	categories, err := sr.selectedUC.GetAllSelectedCategoriesByUser(userId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"selected_categories": categories,
	})
}

func (sr *SelectedRouter) GetUsersCountToSelectedModuleOrCategory(c echo.Context) error {
	isModule, isCategory := true, true

	moduleId, err := strconv.Atoi(c.QueryParam("module_id"))
	if err != nil {
		isModule = false
	}
	categoryId, err := strconv.Atoi(c.QueryParam("category_id"))
	if err != nil {
		isCategory = false
	}

	if isCategory == isModule {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request param",
		})
	}

	var count int

	if isModule {
		count, err = sr.selectedUC.GetUsersCountToSelectedModule(moduleId)
	} else if isCategory {
		count, err = sr.selectedUC.GetUsersCountToSelectedCategory(categoryId)
	}

	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

func (sr *SelectedRouter) InsertSelectedModuleToUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	moduleId, err := strconv.Atoi(c.QueryParam("module_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	err = sr.selectedUC.InsertSelectedModuleToUser(userId, moduleId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.NoContent(http.StatusOK)
}

func (sr *SelectedRouter) InsertSelectedCategoryToUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	categoryId, err := strconv.Atoi(c.QueryParam("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	err = sr.selectedUC.InsertSelectedCategoryToUser(userId, categoryId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.NoContent(http.StatusOK)
}

func (sr *SelectedRouter) DeleteModuleToUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	moduleId, err := strconv.Atoi(c.QueryParam("module_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	err = sr.selectedUC.DeleteModuleToUser(userId, moduleId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.NoContent(http.StatusOK)
}

func (sr *SelectedRouter) DeleteCategoryToUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	categoryId, err := strconv.Atoi(c.QueryParam("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	err = sr.selectedUC.DeleteCategoryToUser(userId, categoryId)
	if err != nil {
		return c.JSON(sr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.NoContent(http.StatusOK)
}
