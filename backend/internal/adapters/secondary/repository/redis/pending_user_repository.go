// Package repository_redis provides Redis-based implementations for persistence ports.
package repository_redis

import (
	"context"
	"strconv"
	"time"

	modelsdb "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/database"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/redis/go-redis/v9"
)

// RedisPendingUserRepository stores pending users as Redis Hashes.
type RedisPendingUserRepository struct {
	client  *redis.Client
	context context.Context
}

var _ output.PendingUserRepository = (*RedisPendingUserRepository)(nil)

// NewRedisPendingUserRepository creates a repository backed by Redis.
func NewRedisPendingUserRepository(client *redis.Client) output.PendingUserRepository {
	return &RedisPendingUserRepository{
		client:  client,
		context: context.Background(),
	}
}

// SavePendingUser stores every pending-user field in a Redis Hash and applies its TTL.
func (r *RedisPendingUserRepository) SavePendingUser(user *modelsdb.PendingUser, ttl time.Duration) error {
	key := "pending_user:" + user.ID
	fields := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"email":         user.Email,
		"hash_code":     user.HashCode,
		"attempts":      strconv.Itoa(user.Attempts),
	}

	pipeline := r.client.TxPipeline()
	pipeline.HSet(r.context, key, fields)
	pipeline.Expire(r.context, key, ttl)
	_, err := pipeline.Exec(r.context)
	return err
}
