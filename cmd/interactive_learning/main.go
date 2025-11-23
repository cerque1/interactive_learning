package main

import (
	"embed"
	"interactive_learning/internal/interactive_learning"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const path_to_static = "./static"

func main() {
	interactive_learning.Run(migrationsDir, MigrationsFS, path_to_static)
}
