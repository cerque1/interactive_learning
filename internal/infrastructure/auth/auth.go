package auth

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/usecase"
	"interactive_learning/internal/utils/tokengenerator"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthRoutes struct {
	UsersUC  usecase.Users
	TokensUC usecase.Tokens
}

func NewAuthRoutes(usersUC usecase.Users, tokensUC usecase.Tokens) *AuthRoutes {
	return &AuthRoutes{UsersUC: usersUC, TokensUC: tokensUC}
}

func (auth *AuthRoutes) AuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" || !strings.HasPrefix("token", "Bearer") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "bad token",
			})
		}

		token = strings.TrimPrefix(token, "Bearer ")
		user_id, err := auth.TokensUC.IsValidToken(tokengenerator.Token(token))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "bad token",
			})
		}

		c.QueryParams().Add("user_id", strconv.Itoa(user_id))

		return next(c)
	}
}

func (auth *AuthRoutes) Login(c echo.Context) error {
	login := c.QueryParam("login")
	password := c.QueryParam("password")

	if login == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "wrong data",
		})
	}

	user, err := auth.UsersUC.GetUserByLogin(login)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unvalid login",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "unauthorized",
		})
	}

	token := auth.TokensUC.AddTokenToUser(user.Id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (auth *AuthRoutes) Register(c echo.Context) error {
	login := c.QueryParam("login")
	name := c.QueryParam("name")
	password := c.QueryParam("password")

	if login == "" || name == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "wrong data",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "server error",
		})
	}

	id, err := auth.UsersUC.InsertUser(
		entity.User{
			Login:        login,
			Name:         name,
			PasswordHash: string(hash),
		})

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "server error",
		})
	}

	token := auth.TokensUC.AddTokenToUser(id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    id,
		"token": token,
	})
}
