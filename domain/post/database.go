package post

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgreDatabase() (*sql.DB, error) {
	connStr := "user=pgblog dbname=pgblog sslmode=verify-full"
	// pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
	// 	os.Getenv("PGHOST"),
	// 	os.Getenv("PGPORT"),
	// 	os.Getenv("PGDATABASE"),
	// 	os.Getenv("PGUSER"),
	// 	os.Getenv("PGPASSWORD"),
	// )
	// connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	db, _ := sql.Open("postgres", connStr)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}
