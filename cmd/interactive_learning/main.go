package main

import (
	"embed"
	"interactive_learning/internal/interactive_learning"
)

const MigrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const pathToStatic = "./static"

func main() {
	interactive_learning.Run(MigrationsDir, MigrationsFS, pathToStatic)
}
