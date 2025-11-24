package interactive_learning

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"interactive_learning/internal/infrastructure"
	"interactive_learning/internal/migrator"
	"interactive_learning/internal/repo/persistent"
	interactivelearning "interactive_learning/internal/usecase/interactive_learning"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func Run(migrationsDir string, migrationsFS embed.FS, path_to_static string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_ADDR"),
		os.Getenv("POSTGRES_DB"))

	// first conn for migrations
	db, err := sql.Open(
		"postgres",
		connectionString)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("database not ready")
	}

	migrator := migrator.NewMigrator(migrationsFS, migrationsDir)
	migrator.ApplyMigrations(db)

	// second conn for usecases
	db, err = sql.Open(
		"postgres",
		connectionString)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("database not ready")
	}

	us := interactivelearning.New(persistent.NewUsersRepo(db), persistent.NewCardsRepo(db), persistent.NewModulesRepo(db), persistent.NewCategoryRepo(db), persistent.NewCategoryModulesRepo(db))
	e := infrastructure.NewEcho(path_to_static, us, us, us, us, us, us)

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	go func() {
		log.Println("starting server...")
		e.Start(":8080")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}
