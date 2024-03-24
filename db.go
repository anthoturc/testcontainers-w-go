package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgConf struct {
	Host     string
	Port     string
	UserName string
	Password string
	Database string
	SSLMode  string
}

// DefaultPgConf returns a pointer to a default Postgres configuration.
func DefaultPgConf() *PgConf {
	return &PgConf{
		Host:     "localhost",
		Port:     "5432",
		UserName: "user",
		Password: "admin1",
		Database: "demo",
		SSLMode:  "disable",
	}
}

// String will craft a Postgres connection string (i.e. dataSourceName for sql.Open)
func (pgConf *PgConf) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgConf.Host, pgConf.Port, pgConf.UserName, pgConf.Password, pgConf.Database, pgConf.SSLMode,
	)
}

// Open will return an instance of a sql.DB connection.
// The caller is responsible for calling Close() on the connection.
func Open(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}
