package bootstrap_database

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
)

func SetupDatabaseRedis(appConfig *config.AppConfig) (models.RedisConfig, error) {
	redisConfig := appConfig.GetRedisConfig()
	return redisConfig, nil
}