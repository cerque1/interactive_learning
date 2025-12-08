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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
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
	idStr := c.Param("id")
	var id int
	var err error

	if idStr == "" {
		id, err = strconv.Atoi(c.QueryParam("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad user id",
			})
		}
	} else {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad id",
			})
		}
	}

	isFull, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		isFull = false
	}

	categories, err := cr.CategoriesUC.GetCategoriesToUser(id, isFull)
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	isFull, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		isFull = false
	}

	modules, err := cr.CategoryModulesUC.GetModulesToCategory(id, isFull)
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
	userIdStr := c.QueryParam("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	category.OwnerId = userId
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
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	categoryId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad category id",
		})
	}

	moduleId, err := strconv.Atoi(c.QueryParam("module_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	category, err := cr.CategoriesUC.GetCategoryById(categoryId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if category.OwnerId != userId {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "you are not the owner",
		})
	}
	if idx := slices.IndexFunc(category.Modules, func(elt entity.Module) bool { return elt.Id == moduleId }); idx >= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "module is already exists",
		})
	}

	if err = cr.CategoryModulesUC.InsertModuleToCategory(categoryId, moduleId); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
