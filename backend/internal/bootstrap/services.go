package bootstrap

import (
	repository_mysql "github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/repository/mysql"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_comments"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_products"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/jmoiron/sqlx"
)

func SetupCommentService(db *sqlx.DB) (input.CommentGetService, input.CommentAddService) {
	commentRepo := repository_mysql.NewSqlCommentRepository(db)
	commentValidator := &service_comments.CommentValidator{}
	return service_comments.NewCommentGetService(commentRepo), service_comments.NewCommentAddService(commentRepo, commentValidator)
}

func SetupProductsService(db *sqlx.DB) input.ProductsGetService {
	productsRepo := repository_mysql.NewSqlProductsRepository(db)
	return service_products.NewProductsGetService(productsRepo)
}


func SetupUserService(userRepo output.UserRepository, csrfService input.CSRFService) (input.UserServiceLogin, input.UserServiceRegister) {
	userNameValidator := &service_auth.UserNameValidator{}
	passwordValidator := &service_auth.PasswordValidator{}
	
	return service_auth.NewUserLoginService(userRepo, userNameValidator, passwordValidator, csrfService, nil), service_auth.NewUserRegisterService(userRepo, userNameValidator, passwordValidator, csrfService, nil)
}