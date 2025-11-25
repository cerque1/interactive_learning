package category

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/usecase"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryRoutes struct {
	CategoriesUC      usecase.Categories
	CategoryModulesUC usecase.CategoryModules
}

func NewCategoryRoutes(CategoriesUC usecase.Categories, CategoryModulesUC usecase.CategoryModules) *CategoryRoutes {
	return &CategoryRoutes{CategoriesUC: CategoriesUC, CategoryModulesUC: CategoryModulesUC}
}

func (cr *CategoryRoutes) GetCategoryById(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	categories, err := cr.CategoriesUC.GetCategoryById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"category": categories,
	})
}

func (cr *CategoryRoutes) GetCategoriesToUser(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	is_full, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		is_full = false
	}

	categories, err := cr.CategoriesUC.GetCategoriesToUser(id, is_full)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

func (cr *CategoryRoutes) GetModulesToCategory(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	is_full, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		is_full = false
	}

	modules, err := cr.CategoryModulesUC.GetModulesToCategory(id, is_full)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"modules": modules,
	})
}

func (cr *CategoryRoutes) InsertCategory(c echo.Context) error {
	category := entity.Category{}
	user_id_str := c.QueryParam("user_id")
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	category.OwnerId = user_id
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	id, err := cr.CategoriesUC.InsertCategory(category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_id": id,
	})
}

func (cr *CategoryRoutes) InsertModuleToCategory(c echo.Context) error {
	user_id, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	category_id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad category id",
		})
	}

	module_id, err := strconv.Atoi(c.QueryParam("module_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	category, err := cr.CategoriesUC.GetCategoryById(category_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if category.OwnerId != user_id {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "you are not the owner",
		})
	}
	if idx := slices.IndexFunc(category.Modules, func(elt entity.Module) bool { return elt.Id == module_id }); idx >= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "module is already exists",
		})
	}

	if err = cr.CategoryModulesUC.InsertModuleToCategory(category_id, module_id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
