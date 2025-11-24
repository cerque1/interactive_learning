package infrastructure

import (
	"interactive_learning/internal/infrastructure/auth"
	"interactive_learning/internal/infrastructure/card"
	"interactive_learning/internal/infrastructure/category"
	"interactive_learning/internal/infrastructure/module"
	"interactive_learning/internal/infrastructure/user"
	"interactive_learning/internal/usecase"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func NewEcho(path_to_static string, usersUC usecase.Users, tokensUC usecase.Tokens, cardUC usecase.Cards, modulesUC usecase.Modules, categorieUC usecase.Categories, categoryModulesUC usecase.CategoryModules) *echo.Echo {
	auth_routes := auth.NewAuthRoutes(usersUC, tokensUC)
	users_routes := user.NewUserRoues(usersUC)
	module_routes := module.NewModuleRoutes(modulesUC, cardUC)
	card_routes := card.NewCardRoutes(cardUC)
	categories_routes := category.NewCategoryRoutes(categorieUC, categoryModulesUC)

	e := echo.New()
	e.Static("/static", path_to_static)

	api := e.Group("/api")

	auth_group := api.Group("/auth")
	auth_group.POST("/login", auth_routes.Login)
	auth_group.POST("/register", auth_routes.Register)

	v1 := api.Group("/v1")
	v1.Use(auth_routes.AuthToken)

	users := v1.Group("/user")
	users.GET("/me", users_routes.GetUserInfoById)
	users.GET("/:id", users_routes.GetUserInfoById)

	categories := v1.Group("/category")
	categories.POST("/:category_id/add_module", categories_routes.InsertModuleToCategory)
	categories.GET("/:id/modules", categories_routes.GetModulesToCategory)
	categories.GET("/:id", categories_routes.GetCategoryById)
	categories.GET("/to_user/:id", categories_routes.GetCategoriesToUser)
	categories.POST("/create", categories_routes.InsertCategory)

	modules := v1.Group("/module")
	modules.GET("/:id", module_routes.GetModuleById)
	modules.GET("/to_user/:id", module_routes.GetModulesByUser)
	modules.POST("/create", module_routes.InsertModule)

	cards := v1.Group("/card")
	cards.GET("/:id", card_routes.GetCardById)
	cards.GET("/to_module/:id", card_routes.GetCardsByModule)
	cards.POST("/create", card_routes.InsertCard)

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	return e
}
