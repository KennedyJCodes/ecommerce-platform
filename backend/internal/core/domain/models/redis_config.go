// Package models defines core domain entities for the sale-watches application.
// This file declares the RedisConfig struct, which encapsulates all parameters required to establish and tune a connection to a Redis server.
package models

import "time"

// RedisConfig holds the connection and performance settings for the Redis client.
// Beyond basic credentials, it includes timeout and pooling configurations to ensure the application remains resilient under high load and handles network instability gracefully.
type RedisConfig struct {
	// Host: the network address of the Redis server.
	Host            string

	// Port: the port number on which Redis is listening (default is usually 6379).
	Port            string

	// Username: the identifier for Redis ACL (Access Control List) authentication.
	Username        string

	// Password: the secret key required to access the Redis instance.
	Password        string

	// DB: the specific logical database index to use (default is 0).
	DB              int

	// --- Timeout & Resilience Settings ---

	// DialTimeout: maximum time to wait for a connection to be established.
	DialTimeout     time.Duration

	// ReadTimeout: maximum time to wait for a reply from Redis after a command is sent.
	ReadTimeout     time.Duration

	// WriteTimeout: maximum time to wait for a command to be written to the network.
	WriteTimeout    time.Duration

	// --- Connection Pooling Settings ---

	// PoolSize: the maximum number of simultaneous connections in the pool.
	PoolSize        int

	// MinIdleConns: the minimum number of idle connections to keep in the pool.
	MinIdleConns    int

	// MaxRetries: number of times to retry a failed command before returning an error.
	MaxRetries      int
}