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

func NewEcho(pathToStatic string, usersUC usecase.Users, tokensUC usecase.Tokens, cardUC usecase.Cards, modulesUC usecase.Modules, categorieUC usecase.Categories, categoryModulesUC usecase.CategoryModules) *echo.Echo {
	authRoutes := auth.NewAuthRoutes(usersUC, tokensUC)
	usersRoutes := user.NewUserRoues(usersUC)
	moduleRoutes := module.NewModuleRoutes(modulesUC, cardUC)
	cardRoutes := card.NewCardRoutes(cardUC)
	categoriesRoutes := category.NewCategoryRoutes(categorieUC, categoryModulesUC)

	e := echo.New()
	e.Static("/static", pathToStatic)

	api := e.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.POST("/login", authRoutes.Login)
	authGroup.POST("/register", authRoutes.Register)

	v1 := api.Group("/v1")
	v1.Use(authRoutes.AuthToken)

	users := v1.Group("/user")
	users.GET("/me", usersRoutes.GetUserInfoById)
	users.GET("/:id", usersRoutes.GetUserInfoById)

	categories := v1.Group("/category")
	categories.POST("/:category_id/add_module", categoriesRoutes.InsertModuleToCategory)
	categories.GET("/:id/modules", categoriesRoutes.GetModulesToCategory)
	categories.GET("/:id", categoriesRoutes.GetCategoryById)
	categories.GET("/to_user/:id", categoriesRoutes.GetCategoriesToUser)
	categories.POST("/create", categoriesRoutes.InsertCategory)

	modules := v1.Group("/module")
	modules.GET("/:id", moduleRoutes.GetModuleById)
	modules.GET("/to_user/:id", moduleRoutes.GetModulesByUser)
	modules.POST("/create", moduleRoutes.InsertModule)

	cards := v1.Group("/card")
	cards.GET("/:id", cardRoutes.GetCardById)
	cards.GET("/to_module/:id", cardRoutes.GetCardsByModule)
	cards.POST("/create", cardRoutes.InsertCard)

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	return e
}
