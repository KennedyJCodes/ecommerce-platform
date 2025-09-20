package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/ports"
)

// DomainRedisTokenRepository implements the domain TokenRepository interface using Redis.
type DomainRedisTokenRepository struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

// NewDomainRedisTokenRepository creates a new Redis token repository instance using domain ports.
func NewDomainRedisTokenRepository(addr string) ports.TokenRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &DomainRedisTokenRepository{
		client: rdb,
		ctx:    context.Background(),
		key:    "paypal_access_token",
	}
}

// GetValidToken retrieves a valid token from Redis storage.
func (r *DomainRedisTokenRepository) GetValidToken() (string, bool, error) {
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