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
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	is_with_cards, err := strconv.ParseBool(c.QueryParam("with_cards"))
	if err != nil {
		log.Println(err)
		is_with_cards = false
	}

	var modules []entity.Module
	if !is_with_cards {
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
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
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
	user_id_str := c.QueryParam("user_id")
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	module.OwnerId = user_id
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
