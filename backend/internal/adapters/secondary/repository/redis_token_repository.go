package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
)

// RedisTokenRepository implements TokenRepository using Redis as storage.
type RedisTokenRepository struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

// NewRedisTokenRepository creates a new Redis token repository instance.
func NewRedisTokenRepository(addr string) output.TokenRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisTokenRepository{
		client: rdb,
		ctx:    context.Background(),
		key:    "paypal_access_token",
	}
}

// GetValidToken retrieves a valid token from Redis storage.
func (r *RedisTokenRepository) GetValidToken() (string, bool, error) {
	// 1. Verificar si el token existe
	token, err := r.client.Get(r.ctx, r.key).Result()
	if err == redis.Nil {
		// No existe
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}

	// 2. Verificar TTL (tiempo de vida)
	ttl, err := r.client.TTL(r.ctx, r.key).Result()
	if err != nil {
		return "", false, err
	}

	if ttl > 0*time.Second {
		// Existe y es válido
		return token, true, nil
	}

	// Existe pero ya expiró
	return "", false, nil
}