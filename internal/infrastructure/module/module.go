package module

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ModuleRoutes struct {
	ModuleUC usecase.Modules
	CardUC   usecase.Cards
}

func NewModuleRoutes(moduleUC usecase.Modules, cardUC usecase.Cards) *ModuleRoutes {
	return &ModuleRoutes{ModuleUC: moduleUC, CardUC: cardUC}
}

func (mr *ModuleRoutes) GetModulesByUser(c echo.Context) error {
	idStr := c.Param("id")
	var id int
	var err error

	if idStr == "" {
		id, err = strconv.Atoi(c.QueryParam("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad id " + idStr,
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

	isWithCards, err := strconv.ParseBool(c.QueryParam("with_cards"))
	if err != nil {
		log.Println(err)
		isWithCards = false
	}

	var modules []entity.Module
	if !isWithCards {
		modules, err = mr.ModuleUC.GetModulesByUser(id)
	} else {
		modules, err = mr.ModuleUC.GetModulesWithCardsByUser(id)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"modules": modules,
	})
}

func (mr *ModuleRoutes) GetModuleById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	module, err := mr.ModuleUC.GetModuleById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"module": module,
	})
}

func (mr *ModuleRoutes) InsertModule(c echo.Context) error {
	module := entity.Module{}
	userIdStr := c.QueryParam("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	module.OwnerId = userId
	if err := c.Bind(&module); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	id, ids, err := mr.ModuleUC.InsertModule(module)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_module_id": id,
		"new_cards_ids": ids,
	})
}
