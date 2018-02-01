package db

import (
	"database/sql"
	"fmt"
)

func NewPostgresHandler(opts DatabaseConnectionOpts) (*sql.DB, error) {
	postgresOpts := NewPostgresConnectionOpts(opts)
	connectionString := postgresOpts.GetConnString()

	db, err := sql.Open(Postgres.String(), connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(opts.MaxConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}

type PostgresConnectionOpts struct {
	DatabaseConnectionOpts
	SSLEnabled bool
}

func NewPostgresConnectionOpts(opts DatabaseConnectionOpts) *PostgresConnectionOpts {
	return &PostgresConnectionOpts{
		DatabaseConnectionOpts: opts,
		SSLEnabled:             false,
	}
}

func (opts *PostgresConnectionOpts) GetConnString() string {
	if opts.Port == "" {
		opts.Port = "3306"
	}
	if opts.Timeout < 0 {
		opts.Timeout = 0
	}
	sslMode := "disable"

	if opts.SSLEnabled {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"host=%v dbname=%v port=%v user=%v password='%v' connect_timeout=%d sslmode=%v",
		opts.Host,
		opts.Database,
		opts.Port,
		opts.User,
		opts.Password,
		opts.Timeout,
		sslMode,
	)
}
