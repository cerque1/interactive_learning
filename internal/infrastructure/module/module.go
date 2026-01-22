package module

import (
	"interactive_learning/internal/entity"
	httputils "interactive_learning/internal/http_utils"
	errors_mapper "interactive_learning/internal/mappers/errors"
	"interactive_learning/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ModuleRoutes struct {
	ModuleUC usecase.Modules
	CardUC   usecase.Cards

	errorsMapper *errors_mapper.ApplicationErrorsMapper
}

func NewModuleRoutes(moduleUC usecase.Modules, cardUC usecase.Cards, errorsMapper *errors_mapper.ApplicationErrorsMapper) *ModuleRoutes {
	return &ModuleRoutes{ModuleUC: moduleUC, CardUC: cardUC, errorsMapper: errorsMapper}
}

func (mr *ModuleRoutes) GetModulesByUser(c echo.Context) error {
	idStr := c.Param("id")
	var id int
	var err error

	id, err = strconv.Atoi(idStr)
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

	isWithCards, err := strconv.ParseBool(c.QueryParam("with_cards"))
	if err != nil {
		isWithCards = false
	}

	modules, err := mr.ModuleUC.GetModulesByUser(id, isWithCards, userId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
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

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	module, err := mr.ModuleUC.GetModuleById(id, userId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"module": module,
	})
}

func (mr *ModuleRoutes) GetModulesByIds(c echo.Context) error {
	modulesIds := httputils.GetModulesByIdsReq{}
	if err := c.Bind(&modulesIds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad data",
		})
	}

	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	isWithCards, err := strconv.ParseBool(c.QueryParam("with_cards"))
	if err != nil {
		isWithCards = false
	}

	modules, err := mr.ModuleUC.GetModulesByIds(modulesIds.ModulesIds, isWithCards, userId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"modules": modules,
	})
}

func (mr *ModuleRoutes) GetPopularModule(c echo.Context) error {
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

	popularModules, err := mr.ModuleUC.GetPopularModules(limit, offset)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"popular_modules": popularModules,
	})
}

func (mr *ModuleRoutes) SearchModules(c echo.Context) error {
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

	foundModules, err := mr.ModuleUC.GetModulesWithSimilarName(name, limit, offset, userId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"found_modules": foundModules,
	})
}

func (mr *ModuleRoutes) InsertModule(c echo.Context) error {
	moduleReq := httputils.ModuleCreateReq{}
	module := entity.ModuleToCreate{}
	userIdStr := c.QueryParam("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	if err := c.Bind(&moduleReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	module.OwnerId = userId
	module.Name = moduleReq.Name
	module.Type = moduleReq.Type

	id, ids, err := mr.ModuleUC.InsertModule(module)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_module_id": id,
		"new_cards_ids": ids,
	})
}

func (mr *ModuleRoutes) RenameModule(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	idStr := c.Param("id")
	moduleId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	newName := httputils.RenameReq{}
	if err = c.Bind(&newName); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad data",
		})
	}

	err = mr.ModuleUC.RenameModule(userId, moduleId, newName.NewName)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (mr *ModuleRoutes) ChangeModuleType(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}
	moduleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module id",
		})
	}
	var moduleType httputils.TypeFromReq
	if err := c.Bind(&moduleType); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad module type",
		})
	}

	err = mr.ModuleUC.UpdateModuleType(moduleId, moduleType.Type, userId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.NoContent(http.StatusOK)
}

func (mr *ModuleRoutes) DeleteModule(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	idStr := c.Param("id")
	moduleId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	err = mr.ModuleUC.DeleteModule(userId, moduleId)
	if err != nil {
		return c.JSON(mr.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
