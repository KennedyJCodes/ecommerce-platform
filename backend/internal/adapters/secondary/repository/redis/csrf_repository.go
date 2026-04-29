// Package repository_redis provides Redis-based implementations for high-performance data persistence.
// This file contains the RedisCSRFRepository, which manages the lifecycle of CSRF tokens, ensuring fast access and automatic expiration via Redis TTL.
package repository_redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// RedisCSRFRepository implements a repository for CSRF tokens using Redis as the backend store.
// It leverages Redis key-value storage to manage security tokens associated with user sessions.
type RedisCSRFRepository struct {
	// client: the Redis client used to execute commands.
	client *redis.Client
	// context: the background context for Redis operations.
	context    context.Context
}

// NewRedisCSRFRepository creates and returns a new instance of RedisCSRFRepository.
// Parameters:
//   - client: an initialized *redis.Client connection.
func NewRedisCSRFRepository(client *redis.Client) *RedisCSRFRepository {
	return &RedisCSRFRepository{
		client: client,
		context:    context.Background(),
	}
}

// Save serializes and persists a CSRF token in Redis.
// It calculates the Time-To-Live (TTL) based on the token's expiration date.
// The key is prefixed with "csrf:" followed by the UserID for easy identification.
// Parameters:
//   - token: a pointer to the models.CSRFToken containing UserID, Token value, and Expiration.
// Returns:
//   - error: an error if serialization or the Redis SET command fails.
func (r *RedisCSRFRepository) Save(token *models.CSRFToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	key := "csrf:" + token.UserID
	// Calculate TTL dynamically based on the distance to ExpiresAt.
	timeToLive := time.Until(token.ExpiresAt)
	
	return r.client.Set(r.context, key, data, timeToLive).Err()
}

// Find retrieves a CSRF token from Redis for a specific user.
// If the token has expired or does not exist, it returns a NotFoundError.
// Otherwise, it deserializes the JSON data into a models.CSRFToken struct.
// Parameters:
//   - userID: the unique identifier of the user whose token is being searched.
// Returns:
//   - *models.CSRFToken: the retrieved token if found.
//   - error: a NotFoundError if missing, or InternalError for other Redis/parsing failures.
func (r *RedisCSRFRepository) Find(userID string) (*models.CSRFToken, error) {
	key := "csrf:" + userID
	data, err := r.client.Get(r.context, key).Bytes()
	if err == redis.Nil {
		// Specific mapping for expired or non-existent keys in Redis.
		return nil, errors.NewNotFoundError(errors.ErrCSRFTokenNotFound)
	}
	if err != nil {
		return nil, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}

	var token models.CSRFToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}

	return &token, nil	
}

// Delete removes a CSRF token from Redis, effectively invalidating the session's CSRF state.
// Parameters:
//   - userID: the unique identifier of the user whose token should be removed.
// Returns:
//   - error: an error if the Redis DEL command fails.
func (r *RedisCSRFRepository) Delete(userID string) error {
	key := "csrf:" + userID
	return r.client.Del(r.context, key).Err()
}