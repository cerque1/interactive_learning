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

	user, err := ur.UsersUC.GetUserInfoById(id, isFull)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
