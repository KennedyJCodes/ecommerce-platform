// Package bootstrap_database provides functions to initialize and configure database connections.
// This file specifically handles the initialization of the Redis client.
package bootstrap_database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KennedyJCodes/ecommerce-platform/internal/config"
	"github.com/redis/go-redis/v9"
)

// SetupRedis initializes and validates a Redis client using the provided configuration.
// It returns an error when Redis cannot be reached, allowing the application to decide how to handle startup failure.
func SetupRedis(cfg config.RedisConfig) (*redis.Client, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("error connecting to Redis on %s: %w", addr, err)
	}

	log.Printf("Successfully connected to Redis on %s - Response: %s", addr, pong)
	return client, nil
}
