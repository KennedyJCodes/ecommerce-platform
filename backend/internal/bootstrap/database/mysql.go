package bootstrap_database

import (
	"fmt"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/jmoiron/sqlx"
)

func SetupDatabaseMySQL(appConfig *config.AppConfig) (*sqlx.DB, error) {
	cfg := appConfig.GetConfig()
	user := cfg.GetString("database.user")
	password := cfg.GetString("database.password")
	host := cfg.GetString("database.host")
	port := cfg.GetInt("database.port")
	dbName := cfg.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)
	return sqlx.Connect("mysql", dsn)
}