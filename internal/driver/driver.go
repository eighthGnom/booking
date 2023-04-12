package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleConn = 5
const maxDBLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = db
	return dbConn, nil
}

// NewDatabase creates new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
