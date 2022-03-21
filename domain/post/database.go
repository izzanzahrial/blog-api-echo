package post

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgreDatabase() (*sql.DB, error) {
	connStr := "user=pgblog dbname=pgblog sslmode=verify-full"
	// connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	db, _ := sql.Open("postgres", connStr)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}
