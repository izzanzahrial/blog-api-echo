package postgre

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgreDatabase() (*sql.DB, error) {
	// another way to create postgresql connection
	// connStr := "user=pgblog dbname=pgblog sslmode=verify-full"
	// connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)
	db, _ := sql.Open("postgres", pgConnString)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}
