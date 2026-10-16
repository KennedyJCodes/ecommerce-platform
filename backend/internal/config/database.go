package config

// IsDatabaseTLSEnabled returns whether the MySQL connection should use TLS.
func (a *AppConfig) IsDatabaseTLSEnabled() bool {
	return a.config.GetBool("DATABASE_TLS")
}
