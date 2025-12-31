package repository_redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisCSRFRepository struct {
	client *redis.Client
	context    context.Context
}

func NewRedisCSRFRepository(client *redis.Client) *RedisCSRFRepository {
	return &RedisCSRFRepository{
		client: client,
		context:    context.Background(),
	}
}

func (r *RedisCSRFRepository) Save(token *models.CSRFToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	key := "csrf:" + token.UserID
	timeToLive := time.Until(token.ExpiresAt)
	
	return r.client.Set(r.context, key, data, timeToLive).Err()
}

func (r *RedisCSRFRepository) Find(userID string) (*models.CSRFToken, error) {
	key := "csrf:" + userID
	data, err := r.client.Get(r.context, key).Bytes()
	if err == redis.Nil {
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

func (r *RedisCSRFRepository) Delete(userID string) error {
	key := "csrf:" + userID
	return r.client.Del(r.context, key).Err()
}