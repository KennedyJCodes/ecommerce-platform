package config

// GetJWTSecret retrieves the JWT secret key from configuration.
func (a *AppConfig) GetJWTSecret() string {
	return a.config.GetString("SECURITY_JWT_JWT_SECRET")
}

// GetCookiePrefix returns the prefix prepended to all cookie names.
// Used to avoid collisions when multiple applications share a parent domain.
func (a *AppConfig) GetCookiePrefix() string {
	return a.config.GetString("COOKIE_PREFIX")
}

// IsSSLEnabled returns whether SSL/TLS should be used.
func (a *AppConfig) IsSSLEnabled() bool {
	return a.config.GetBool("SSL_ENABLED")
}

// GetSSLCertFile returns the path to the SSL certificate file.
func (a *AppConfig) GetSSLCertFile() string {
	return a.config.GetString("SSL_CERT_FILE")
}

// GetSSLKeyFile returns the path to the SSL private key file.
func (a *AppConfig) GetSSLKeyFile() string {
	return a.config.GetString("SSL_KEY_FILE")
}
