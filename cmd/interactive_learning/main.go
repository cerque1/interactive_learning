package main

import (
	"embed"
	"interactive_learning/internal/interactive_learning"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	interactive_learning.Run(migrationsDir, MigrationsFS)
}
