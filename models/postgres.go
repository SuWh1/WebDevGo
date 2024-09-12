package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Open will opens SQL connection
// if we call Open() we should to ensure that it will be closed with db.Close() method
func Open(config PostgresConfig) (*sql.DB, error) { // database connection
	db, err := sql.Open("pgx", config.String()) // format PostgreConfig strings into a single connection string
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "dauren",
		Password: "daukahifi",
		Database: "suwh",
		SSLMode:  "disable",
	} // just for local use
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// pgx is driver name
// open creates connection, but not make request
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
