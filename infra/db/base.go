package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"gopkg.in/doug-martin/goqu.v4"
)

type Engine string

const (
	MySQL    Engine = "mysql"
	Postgres Engine = "postgres"
)

func (e Engine) String() string {
	return string(e)
}

func New(engine Engine, opts DatabaseConnectionOpts) (*goqu.Database, error) {
	var (
		db  *sql.DB
		err error
	)
	switch engine {
	case MySQL:
		db, err = NewMySQLHandler(opts)
	case Postgres:
		fallthrough
	default:
		engine = Postgres
		db, err = NewMySQLHandler(opts)
	}
	if err != nil {
		return nil, err
	}
	return goqu.New(engine.String(), db), nil
}

type DatabaseConnectionOpts struct {
	Host               string
	Database           string
	User               string
	Password           string
	Port               string
	Timeout            int
	MaxConnections     int
	MaxIdleConnections int
	Params             map[string]string
}

func NewDatabaseConnectionOpts(host, database, user, password, port string, timeout, maxConnections, maxIdleConnections int) *DatabaseConnectionOpts {
	return &DatabaseConnectionOpts{
		Host:               host,
		Database:           database,
		User:               user,
		Password:           password,
		Port:               port,
		Timeout:            timeout,
		MaxConnections:     maxConnections,
		MaxIdleConnections: maxIdleConnections,
	}
}

func (opts *DatabaseConnectionOpts) GetConnString() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s/%s?%s",
		opts.User,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		opts.GetParams,
	)
}

func (opts *DatabaseConnectionOpts) GetParams() string {
	if opts.Params == nil || len(opts.Params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for k, v := range opts.Params {
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)
		buffer.WriteString("&")
	}
	str := buffer.String()
	return str[:len(str)-1]
}
