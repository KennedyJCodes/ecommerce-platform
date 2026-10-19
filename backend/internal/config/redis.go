package config

// GetRedisConfig returns a RedisConfig struct populated from the REDIS_* configuration keys.
func (a *AppConfig) GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         a.config.GetString("REDIS_HOST"),
		Port:         a.config.GetString("REDIS_PORT"),
		Username:     a.config.GetString("REDIS_USERNAME"),
		Password:     a.config.GetString("REDIS_PASSWORD"),
		DB:           a.config.GetInt("REDIS_DB"),
		DialTimeout:  a.config.GetDuration("REDIS_DIAL_TIMEOUT"),
		ReadTimeout:  a.config.GetDuration("REDIS_READ_TIMEOUT"),
		WriteTimeout: a.config.GetDuration("REDIS_WRITE_TIMEOUT"),
		PoolSize:     a.config.GetInt("REDIS_POOL_SIZE"),
		MinIdleConns: a.config.GetInt("REDIS_MIN_IDLE_CONNS"),
		MaxRetries:   a.config.GetInt("REDIS_MAX_RETRIES"),
	}
}
