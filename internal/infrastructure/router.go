package infrastructure

import (
	"interactive_learning/internal/infrastructure/auth"
	"interactive_learning/internal/infrastructure/card"
	"interactive_learning/internal/infrastructure/category"
	"interactive_learning/internal/infrastructure/module"
	"interactive_learning/internal/infrastructure/results"
	"interactive_learning/internal/infrastructure/user"
	"interactive_learning/internal/usecase"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func NewEcho(pathToStatic string, usersUC usecase.Users, tokensUC usecase.Tokens, cardUC usecase.Cards, modulesUC usecase.Modules, categorieUC usecase.Categories, categoryModulesUC usecase.CategoryModules, resultsUC usecase.Results) *echo.Echo {
	authRoutes := auth.NewAuthRoutes(usersUC, tokensUC)
	usersRoutes := user.NewUserRoues(usersUC)
	moduleRoutes := module.NewModuleRoutes(modulesUC, cardUC)
	cardRoutes := card.NewCardRoutes(cardUC)
	categoriesRoutes := category.NewCategoryRoutes(categorieUC, categoryModulesUC)
	resultsRoutes := results.NewResultsRoutes(resultsUC)

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

	results := v1.Group("/results")
	results.GET("/to_user/:id", resultsRoutes.GetResultsByOwner)
	results.GET("/cards_result/:result_id", resultsRoutes.GetCardsResultById)
	results.GET("/category_result/:category_res_id", resultsRoutes.GetCategoryResById)

	moduleResult := results.Group("/module_result")
	moduleResult.POST("/insert", resultsRoutes.InsertModuleResult)
	moduleResult.DELETE("/delete/:id", resultsRoutes.DeleteModuleResult)

	categoryResult := results.Group("category_result")
	categoryResult.POST("/insert", resultsRoutes.InsertCategoryResult)
	categoryResult.DELETE("/delete/:id", resultsRoutes.DeleteCategoryResultById)

	categories := v1.Group("/category")
	categories.POST("/:category_id/add_modules", categoriesRoutes.InsertModulesToCategory)
	categories.DELETE("/:category_id/:module_id/delete", categoriesRoutes.DeleteModuleFromCategory)
	categories.GET("/:id/modules", categoriesRoutes.GetModulesToCategory)
	categories.GET("/:id", categoriesRoutes.GetCategoryById)
	categories.GET("/to_user/:id", categoriesRoutes.GetCategoriesToUser)
	categories.POST("/create", categoriesRoutes.InsertCategory)
	categories.PUT("/rename/:id", categoriesRoutes.RenameCategory)
	categories.DELETE("/delete/:id", categoriesRoutes.DeleteCategory)

	modules := v1.Group("/module")
	modules.GET("/:id", moduleRoutes.GetModuleById)
	modules.GET("/to_user/:id", moduleRoutes.GetModulesByUser)
	modules.POST("/by_ids", moduleRoutes.GetModulesByIds)
	modules.POST("/create", moduleRoutes.InsertModule)
	modules.PUT("/rename/:id", moduleRoutes.RenameModule)
	modules.DELETE("/delete/:id", moduleRoutes.DeleteModule)

	cards := v1.Group("/card")
	cards.GET("/:id", cardRoutes.GetCardById)
	cards.GET("/to_module/:id", cardRoutes.GetCardsByModule)
	cards.POST("/insert_to_module", cardRoutes.InsertCards)
	cards.PUT("/update/:id", cardRoutes.UpdateCard)
	cards.DELETE("/delete/:id", cardRoutes.DeleteCard)

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	return e
}
