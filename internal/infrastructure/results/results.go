package results

import (
	httputils "interactive_learning/internal/http_utils"
	"interactive_learning/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ResultsRoutes struct {
	ResultsUC usecase.Results
}

func NewResultsRoutes(ResultsUC usecase.Results) *ResultsRoutes {
	return &ResultsRoutes{ResultsUC: ResultsUC}
}

func (rr *ResultsRoutes) GetResultsByOwner(c echo.Context) error {
	idStr := c.Param("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	isModuleRes := true
	moduleId, err := strconv.Atoi(c.QueryParam("moduleId"))
	if err != nil {
		isModuleRes = false
	}

	isCategoryRes := true
	categoryId, err := strconv.Atoi(c.QueryParam("categoryId"))
	if err != nil {
		isCategoryRes = false
	}

	if isCategoryRes && isModuleRes {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "The request must contain either a category, a module, or nothing.",
		})
	} else if isModuleRes {
		moduleResults, err := rr.ResultsUC.GetResultsToModuleId(moduleId, userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"module_results": moduleResults,
		})
	} else if isCategoryRes {
		categoryResults, err := rr.ResultsUC.GetResultsByCategoryId(categoryId, userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"category_results": categoryResults,
		})
	}

	categoriesResults, modulesResults, err := rr.ResultsUC.GetResultsByOwner(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories_results": categoriesResults,
		"modules_results":    modulesResults,
	})
}

func (rr *ResultsRoutes) GetModuleResultById(c echo.Context) error {
	idStr := c.Param("result_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad result id",
		})
	}

	moduleRes, err := rr.ResultsUC.GetModuleResultById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"module_res": moduleRes,
	})
}

func (rr *ResultsRoutes) GetCardsResultById(c echo.Context) error {
	idStr := c.Param("result_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad result id",
		})
	}

	cardsResults, err := rr.ResultsUC.GetCardsResultById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cards_results": cardsResults,
	})
}

func (rr *ResultsRoutes) GetCategoryResById(c echo.Context) error {
	idStr := c.Param("category_res_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad result id",
		})
	}

	categoryResult, err := rr.ResultsUC.GetCategoryResById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"category_result": categoryResult,
	})
}

func (rr *ResultsRoutes) InsertModuleResult(c echo.Context) error {
	var insertModuleReq httputils.InsertModuleResultReq
	if err := c.Bind(&insertModuleReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	insertModuleReq.Owner = userId

	newId, err := rr.ResultsUC.InsertModuleResult(insertModuleReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]int{
		"new_id": newId,
	})
}

func (rr *ResultsRoutes) InsertCategoryResult(c echo.Context) error {
	var insertCategoryReq httputils.InsertCategoryModulesResultReq
	if err := c.Bind(&insertCategoryReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	insertCategoryReq.Owner = userId

	newCategoryResultId, newResultsIds, err := rr.ResultsUC.InsertCategoryResult(insertCategoryReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_category_result_id": newCategoryResultId,
		"new_results_ids":        newResultsIds,
	})
}

func (rr *ResultsRoutes) DeleteModuleResult(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad result id",
		})
	}

	err = rr.ResultsUC.DeleteModuleResult(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (rr *ResultsRoutes) DeleteCategoryResultById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad result id",
		})
	}

	err = rr.ResultsUC.DeleteCategoryResultById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}
