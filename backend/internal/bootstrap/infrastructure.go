package bootstrap

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/repository/mysql"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/static"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/jmoiron/sqlx"
)

func SetupStaticFileAdapter(appConfig *config.AppConfig) output.StaticFilePort {
	staticDir := appConfig.GetStaticDir()
	return static.NewStaticFileAdapter(staticDir)
}

func SetupUserRepository(db *sqlx.DB) output.UserRepository {
	hasher := security_auth.BcryptHasher{}
	return repository_mysql.NewSQLUserRepository(db, hasher)
}

func SetupCommonServices(appConfig *config.AppConfig) {
	security_auth.SetDefaultJWTService(appConfig.GetJWTSecret())
}