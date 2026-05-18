// Package repository_redis provides Redis-based implementations for high-performance data persistence.
// This file contains the RedisBlacklistRepository, which manages revoked JWT tokens using Redis keys with automatic TTL-based expiration.
package repository_redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisBlacklistRepository implements the TokenBlacklistPort interface using Redis as the backend store.
// It leverages Redis key-value storage to track revoked JWT tokens, with automatic expiration
// aligned to each token's remaining validity period.
type RedisBlacklistRepository struct {
	// client: the Redis client used to execute commands.
	client *redis.Client
	// context: the background context for Redis operations.
	context context.Context
}

// NewRedisBlacklistRepository creates and returns a new instance of RedisBlacklistRepository.
// Parameters:
//   - client: an initialized *redis.Client connection.
func NewRedisBlacklistRepository(client *redis.Client) *RedisBlacklistRepository {
	return &RedisBlacklistRepository{
		client:  client,
		context: context.Background(),
	}
}

// Add stores a JWT identifier in the blacklist with a TTL matching the token's remaining validity.
// The key is prefixed with "blacklist:" followed by the JTI for easy identification and cleanup.
// Parameters:
//   - jti: the unique JWT ID to blacklist.
//   - ttl: the duration for which the token should remain blacklisted.
// Returns:
//   - error: an error if the Redis SET command fails.
func (r *RedisBlacklistRepository) Add(jti string, ttl time.Duration) error {
	key := "blacklist:" + jti
	return r.client.Set(r.context, key, "1", ttl).Err()
}

// IsBlacklisted checks whether a given JWT ID exists in the blacklist.
// Returns false if the key is missing or has expired naturally via Redis TTL.
// Parameters:
//   - jti: the unique JWT ID to check against the blacklist.
// Returns:
//   - bool: true if the token is blacklisted (revoked).
//   - error: an error if the Redis GET command fails.
func (r *RedisBlacklistRepository) IsBlacklisted(jti string) (bool, error) {
	key := "blacklist:" + jti
	_, err := r.client.Get(r.context, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
