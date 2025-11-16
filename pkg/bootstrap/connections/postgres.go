package connections

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitPostgres() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("bootstrap: failed to connect to DB: %v", err)
	}

	for i := range 10 {
		err = db.PingContext(context.Background())
		if err == nil {
			break
		}

		log.Printf("bootstrap: retrying DB connection (%d)...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("bootstrap: could not ping DB: %v", err)
	}

	return db
}
