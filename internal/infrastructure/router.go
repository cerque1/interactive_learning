package infrastructure

import (
	"interactive_learning/internal/infrastructure/auth"
	"interactive_learning/internal/infrastructure/card"
	"interactive_learning/internal/infrastructure/module"
	"interactive_learning/internal/usecase"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func NewEcho(usersUC usecase.Users, tokensUC usecase.Tokens, cardUC usecase.Cards, modulesUC usecase.Modules) *echo.Echo {
	auth_routes := auth.NewAuthRoutes(usersUC, tokensUC)
	module_routes := module.NewModuleRoutes(modulesUC, cardUC)
	card_routes := card.NewCardRoutes(cardUC)

	e := echo.New()

	api := e.Group("/api")

	auth_group := api.Group("/auth")
	auth_group.POST("/login", auth_routes.Login)
	auth_group.POST("/register", auth_routes.Register)

	v1 := api.Group("/v1")
	v1.Use(auth_routes.AuthToken)

	modules := v1.Group("/module")
	modules.GET("/:id", module_routes.GetModulesByUser)
	modules.POST("/create", module_routes.InsertModule)

	cards := v1.Group("/card")
	cards.GET("/:id", card_routes.GetCardsByModule)
	cards.POST("/create", card_routes.InsertCard)

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	return e
}
