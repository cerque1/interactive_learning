package auth

import (
	"interactive_learning/internal/entity"
	errors_mapper "interactive_learning/internal/mappers/errors"
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

	errorsMapper *errors_mapper.ApplicationErrorsMapper
}

func NewAuthRoutes(usersUC usecase.Users, tokensUC usecase.Tokens, errorsMapper *errors_mapper.ApplicationErrorsMapper) *AuthRoutes {
	return &AuthRoutes{UsersUC: usersUC, TokensUC: tokensUC, errorsMapper: errorsMapper}
}

func (auth *AuthRoutes) AuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer") {
			log.Println("prefix " + token)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "bad token",
			})
		}

		token = strings.TrimPrefix(token, "Bearer ")
		userId, err := auth.TokensUC.IsValidToken(tokengenerator.Token(token))
		if err != nil {
			log.Println(token)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "bad token",
			})
		}

		c.QueryParams().Add("user_id", strconv.Itoa(userId))

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
		return c.JSON(auth.errorsMapper.ApplicationErrorToHttp(err))
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
	isContains, err := auth.UsersUC.IsContainsLogin(login)
	if err != nil {
		return c.JSON(auth.errorsMapper.ApplicationErrorToHttp(err))
	} else if isContains {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the login already exists",
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
		return c.JSON(auth.errorsMapper.ApplicationErrorToHttp(err))
	}

	token := auth.TokensUC.AddTokenToUser(id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    id,
		"token": token,
	})
}
