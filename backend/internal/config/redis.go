package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and returns a new Redis client connected to the configured Redis instance.
// It verifies the connection with a 5-second timeout and terminates the application on failure.
func NewRedisClient(cfg RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	})

	// Verify the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis on %s: %v", addr, err)
	}

	log.Printf("Successfully connected to Redis on %s - Response: %s", addr, pong)
	return client
}

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
