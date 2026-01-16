package user

import (
	"errors"
	myerrors "interactive_learning/internal/errors"
	"interactive_learning/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserRoutes struct {
	UsersUC usecase.Users
}

func NewUserRoues(UsersUC usecase.Users) *UserRoutes {
	return &UserRoutes{UsersUC: UsersUC}
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"found_users": foundUsers,
	})
}
