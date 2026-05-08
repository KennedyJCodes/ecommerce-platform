// Package bootstrap_database provides functions to initialize and configure database connections.
// It abstracts the creation of connection strings and the establishment of pool connections using application configuration.
package bootstrap_database

import (
	"fmt"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/config"
	"github.com/jmoiron/sqlx"
)

// SetupDatabaseMySQL initializes a connection to a MySQL database using the provided application configuration.
// It retrieves database credentials (user, password, host, port, name) from the config provider, constructs a Data Source Name (DSN), and establishes a connection pool using sqlx.

// Implementation Details:
//   - Uses "parseTime=true" in the DSN to ensure MySQL DATE/DATETIME fields are correctly mapped to time.Time in Go.
//   - Returns a thread-safe *sqlx.DB connection pool.
//
// Parameters:
//   - appConfig: a pointer to the application configuration manager.
//
// Returns:
//   - *sqlx.DB: an active database connection pool if successful.
//   - error: an error if the connection cannot be established or configuration is missing.
func SetupDatabaseMySQL(appConfig *config.AppConfig) (*sqlx.DB, error) {
	cfg := appConfig.GetConfig()

	// Retrieve connection parameters from the configuration source.
	user := cfg.GetString("database.user")
	password := cfg.GetString("database.password")
	host := cfg.GetString("database.host")
	port := cfg.GetInt("database.port")
	dbName := cfg.GetString("database.name")

	// Construct the MySQL Data Source Name (DSN).
	// parseTime=true is critical for handling time-based domain models.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)

	// Connect and ping the database to ensure it's reachable.
	return sqlx.Connect("mysql", dsn)
}
