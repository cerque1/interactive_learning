package category

import (
	"errors"
	"interactive_learning/internal/entity"
	myerrors "interactive_learning/internal/errors"
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

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	categories, err := cr.CategoriesUC.GetCategoryById(id, userId)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrNotAvailable):
			return c.JSON(http.StatusNotAcceptable, map[string]string{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"category": categories,
	})
}

func (cr *CategoryRoutes) GetCategoriesToUser(c echo.Context) error {
	idStr := c.Param("id")
	var id int
	var err error

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	if idStr == "" {
		id = userId
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

	categories, err := cr.CategoriesUC.GetCategoriesToUser(id, isFull, userId)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrNotAvailable):
			return c.JSON(http.StatusNotAcceptable, map[string]string{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
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

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	isFull, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		isFull = false
	}

	modules, err := cr.CategoryModulesUC.GetModulesToCategory(id, isFull, userId)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrNotAvailable):
			return c.JSON(http.StatusNotAcceptable, map[string]string{
				"message": err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"modules": modules,
	})
}

func (cr *CategoryRoutes) GetPopularCategories(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error parse limit",
		})
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error parse offset",
		})
	}

	popularCategories, err := cr.CategoriesUC.GetPopularCategories(limit, offset)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"popular_categories": popularCategories,
	})
}

func (cr *CategoryRoutes) SearchCategories(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" || len([]byte(name)) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty or short name",
		})
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	foundCategories, err := cr.CategoriesUC.GetCategoriesWithSimilarName(name, limit, offset, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"found_categories": foundCategories,
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

func (cr *CategoryRoutes) ChangeCategoryType(c echo.Context) error {
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
	var categoryType httputils.TypeFromReq
	if err := c.Bind(&categoryType); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad category type",
		})
	}

	err = cr.CategoriesUC.UpdateCategoryType(categoryId, categoryType.Type, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
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
