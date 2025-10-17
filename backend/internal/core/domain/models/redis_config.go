package models

import "time"

type RedisConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DB              int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	MinIdleConns    int
	MaxRetries      int
}