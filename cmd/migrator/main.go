package main

import (
	"log"

	"github.com/IlyaAGL/avito_autumn_2025/pkg/bootstrap/connections"
	"github.com/IlyaAGL/avito_autumn_2025/pkg/bootstrap/migrations"
)

func main() {
	db_pg := connections.InitPostgres()
	defer func() {
		log.Fatal(db_pg.Close())
	}()

	migrations.RunMigrationsPG(db_pg)
}
