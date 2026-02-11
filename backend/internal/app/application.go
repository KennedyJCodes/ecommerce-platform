package app

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/bootstrap"
	bootstrap_database "github.com/David-Alejandro-Jimenez/sale-watches/internal/bootstrap/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	config      *config.AppConfig
	db          *sqlx.DB
	redisClient *redis.Client
	httpServer  *http.Server	
}

func NewConfigApplication() *Application {
    return &Application{}
}

func (a *Application) LoadConfig() error {
    a.config = config.NewAppConfig()
    config.ValidateRequiredConfig(a.config.GetConfig())
    return nil
}

func (a *Application) SetupCommonServices() error {
    bootstrap.SetupCommonServices(a.config)
    return nil
}

func (a *Application) SetupDatabase() error {
    db, err := bootstrap_database.SetupDatabaseMySQL(a.config)
    if err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }
    a.db = db
    return nil
}

func (a *Application) SetupRedis() error {
    redisConfig, err := bootstrap_database.SetupDatabaseRedis(a.config)
    if err != nil {
        return fmt.Errorf("error connecting to Redis: %w", err)
    }
    a.redisClient = config.NewRedisClient(redisConfig)
    return nil
}

func (a *Application) GetPort() string {
    return a.config.GetPort()
}

func (a *Application) Close() {
    if a.db != nil {
        a.db.Close()
    }
    if a.redisClient != nil {
        a.redisClient.Close()
    }
}