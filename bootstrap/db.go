package bootstrap

import (
	"github.com/spf13/viper"
	"gopkg.in/doug-martin/goqu.v4"
	"monkey/infra/db"
)

func GetDBConfig() db.DatabaseConnectionOpts {
	return *db.NewDatabaseConnectionOpts(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_PORT"),
		viper.GetInt("DB_CONN"),
		viper.GetInt("DB_MAX_CONN"),
		viper.GetInt("DB_MAX_IDLE_CONN"),
	)
}

func GetDBEngine() db.Engine {
	switch viper.GetString("DB_ENGINE") {
	case "mysql":
		return db.MySQL
	default:
		return db.Postgres
	}
}

func GetDB() (*goqu.Database, error) {
	engine, connOpts := GetDBEngine(), GetDBConfig()
	return db.New(engine, connOpts)
}
