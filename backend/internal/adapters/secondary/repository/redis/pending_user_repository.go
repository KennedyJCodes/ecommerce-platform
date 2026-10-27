// Package repository_redis provides Redis-based implementations for persistence ports.
package repository_redis

import (
	"context"
	"strconv"
	"time"

	modelsdb "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/database"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	appErrors "github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// RedisPendingUserRepository stores pending users as Redis Hashes.
type RedisPendingUserRepository struct {
	client *redis.Client
}

var _ output.PendingUserRepository = (*RedisPendingUserRepository)(nil)

// KEYS[1] = pending_user:{id}      ARGV[1] = max attempts
var registerAttemptScript = redis.NewScript(`
	if redis.call('EXISTS', KEYS[1]) == 0 then
  		return -1
	end
	local attempts = redis.call('HINCRBY', KEYS[1], 'attempts', 1)
	if attempts > tonumber(ARGV[1]) then
  		redis.call('DEL', KEYS[1])
  		return -2
	end
	return attempts
`)

// NewRedisPendingUserRepository creates a repository backed by Redis.
func NewRedisPendingUserRepository(client *redis.Client) output.PendingUserRepository {
	return &RedisPendingUserRepository{
		client: client,
	}
}

// GetPendingUser retrieves a pending user from its Redis Hash.
func (r *RedisPendingUserRepository) GetPendingUser(ctx context.Context, userID string) (*modelsdb.PendingUser, error) {
	key := "pending_user:" + userID
	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, appErrors.NewInternalError(appErrors.ErrDatabaseQuery).WithError(err)
	}
	if len(fields) == 0 {
		return nil, appErrors.NewNotFoundError(appErrors.ErrUserNotFound)
	}

	attempts, err := strconv.Atoi(fields["attempts"])
	if err != nil {
		return nil, appErrors.NewInternalError(appErrors.ErrDatabaseQuery).WithError(err)
	}

	return &modelsdb.PendingUser{
		ID:           fields["id"],
		Username:     fields["username"],
		PasswordHash: fields["password_hash"],
		Email:        fields["email"],
		HashCode:     fields["hash_code"],
		Attempts:     attempts,
	}, nil
}

// GetPendingUserHashCode retrieves only the verification-code hash from a pending user Redis Hash.
func (r *RedisPendingUserRepository) GetPendingUserHashCode(ctx context.Context, userID string) (string, error) {
	key := "pending_user:" + userID
	hashCode, err := r.client.HGet(ctx, key, "hash_code").Result()
	if err == redis.Nil {
		return "", appErrors.NewNotFoundError(appErrors.ErrUserNotFound)
	}
	if err != nil {
		return "", appErrors.NewInternalError(appErrors.ErrDatabaseQuery).WithError(err)
	}

	return hashCode, nil
}

func (r *RedisPendingUserRepository) RegisterVerificationAttempt(ctx context.Context, userID string, maxAttempts int) (int, error) {
	key := "pending_user:" + userID
	result, err := registerAttemptScript.Run(ctx, r.client, []string{key}, maxAttempts).Int()
	if err != nil {
		return 0, appErrors.NewInternalError(appErrors.ErrDatabaseUpdate).WithError(err)
	}
	switch result {
	case -1:
		return 0, appErrors.NewNotFoundError(appErrors.ErrUserNotFound)
	case -2:
		return 0, appErrors.NewTooManyRequestsError(appErrors.ErrTooManyVerificationAttempts)
	}
	return result, nil
}

// SavePendingUser stores every pending-user field in a Redis Hash and applies its TTL.
func (r *RedisPendingUserRepository) SavePendingUser(ctx context.Context, user *modelsdb.PendingUser, ttl time.Duration) error {
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
	pipeline.HSet(ctx, key, fields)
	pipeline.Expire(ctx, key, ttl)
	_, err := pipeline.Exec(ctx)
	return err
}

// DeletePendingUser removes a pending user from Redis after verification.
func (r *RedisPendingUserRepository) DeletePendingUser(ctx context.Context, userID string) error {
	key := "pending_user:" + userID
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return appErrors.NewInternalError(appErrors.ErrDatabaseDelete).WithError(err)
	}
	return nil
}
