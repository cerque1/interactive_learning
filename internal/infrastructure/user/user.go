package user

import (
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
	id_str := c.Param("id")
	var id int
	var err error

	if id_str == "" {
		id, err = strconv.Atoi(c.QueryParam("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad user id",
			})
		}
	} else {
		id, err = strconv.Atoi(id_str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "bad id",
			})
		}
	}

	is_full, err := strconv.ParseBool(c.QueryParam("is_full"))
	if err != nil {
		is_full = false
	}

	user, err := ur.UsersUC.GetUserInfoById(id, is_full)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
