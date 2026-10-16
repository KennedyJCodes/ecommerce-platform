package config

// GetResendEmailAPIKey returns the Resend API key from configuration.
func (a *AppConfig) GetResendEmailAPIKey() string {
	return a.config.GetString("RESEND_EMAIL_API_KEY")
}

// GetResendEmailFrom returns the sender email address from configuration.
func (a *AppConfig) GetResendEmailFrom() string {
	return a.config.GetString("RESEND_EMAIL_FROM")
}
