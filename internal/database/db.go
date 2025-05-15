// internal/database/db.go
package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDatabase() (err error) {
	dsn := "postgresql://dia_intern_prod:gtYXUPuU5c32e7l@10.101.0.241:5432/techupdb"
	DB, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func GetDB() *pgxpool.Pool {
	return DB
}
