package category

import (
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/usecase"
	"net/http"
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
	categoryToCreate := entity.CategoryToCreate{}
	userIdStr := c.QueryParam("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	if err := c.Bind(&categoryToCreate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	categoryToCreate.OwnerId = userId

	id, err := cr.CategoriesUC.InsertCategory(categoryToCreate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_id": id,
	})
}

func (cr *CategoryRoutes) InsertModulesToCategory(c echo.Context) error {
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

	modulesIds := httputils.AddModulesToCategoryReq{}
	if err = c.Bind(&modulesIds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad data",
		})
	}

	if err = cr.CategoryModulesUC.InsertModulesToCategory(userId, categoryId, modulesIds.ModulesIds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (cr *CategoryRoutes) RenameCategory(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad category id",
		})
	}

	newName := httputils.RenameReq{}
	if err = c.Bind(&newName); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad data",
		})
	}

	err = cr.CategoriesUC.RenameCategory(userId, categoryId, newName.NewName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (cr *CategoryRoutes) DeleteCategory(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	idStr := c.Param("id")
	cardId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	err = cr.CategoriesUC.DeleteCategory(userId, cardId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "delete category error " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (cr *CategoryRoutes) DeleteModuleFromCategory(c echo.Context) error {
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

	moduleId, err := strconv.Atoi(c.Param("module_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}

	err = cr.CategoryModulesUC.DeleteModuleFromCategory(userId, categoryId, moduleId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error delete module from category",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
