package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLHandler(opts DatabaseConnectionOpts) (*sql.DB, error) {
	mySQLOpts := NewMySQLConnectionOpts(opts)
	connectionString := mySQLOpts.GetConnString()

	db, err := sql.Open(MySQL.String(), connectionString)
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

type MySQLConnectionOpts struct {
	DatabaseConnectionOpts
	Protocol string
}

func NewMySQLConnectionOpts(opts DatabaseConnectionOpts) *MySQLConnectionOpts {
	return &MySQLConnectionOpts{
		DatabaseConnectionOpts: opts,
		Protocol:               "tcp",
	}
}

func (opts *MySQLConnectionOpts) GetConnString() string {
	if opts.Params == nil {
		opts.Params = make(map[string]string)
	}
	if opts.Port == "" {
		opts.Port = "3306"
	}
	if opts.Protocol == "" {
		opts.Protocol = "tcp"
	}

	opts.Params["parseTime"] = "true"

	return fmt.Sprintf(
		"%s:%s@%s/%s?%s",
		opts.User,
		opts.Password,
		opts.getFullHost(),
		opts.Database,
		opts.GetParams(),
	)
}

func (opts *MySQLConnectionOpts) getFullHost() string {
	return fmt.Sprintf(
		// protocol(host:port)
		"%s(%s:%s)",
		opts.Protocol,
		opts.Host,
		opts.Port,
	)
}
