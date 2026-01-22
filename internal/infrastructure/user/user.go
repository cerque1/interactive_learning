package user

import (
	errors_mapper "interactive_learning/internal/mappers/errors"
	"interactive_learning/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserRoutes struct {
	UsersUC usecase.Users

	errorsMapper *errors_mapper.ApplicationErrorsMapper
}

func NewUserRoues(UsersUC usecase.Users, errorsMapper *errors_mapper.ApplicationErrorsMapper) *UserRoutes {
	return &UserRoutes{UsersUC: UsersUC, errorsMapper: errorsMapper}
}

func (ur *UserRoutes) GetUserInfoById(c echo.Context) error {
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
				"message": "bad user id",
			})
		}
	}

	isFull, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		isFull = false
	}

	user, err := ur.UsersUC.GetUserInfoById(id, isFull, userId)
	if err != nil {
		return c.JSON(ur.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (ur *UserRoutes) SearchUsers(c echo.Context) error {
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

	foundUsers, err := ur.UsersUC.GetUsersWithSimilarName(name, limit, offset)
	if err != nil {
		return c.JSON(ur.errorsMapper.ApplicationErrorToHttp(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"found_users": foundUsers,
	})
}
