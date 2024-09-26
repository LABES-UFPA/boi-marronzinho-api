package usuario

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var UserModule = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewUserUseCase),
	fx.Provide(NewUserHandler),
)

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

func NewUserUseCase(repo repository.UserRepository) *usecase.UsuarioUseCase {
	return usecase.NewUsuarioUseCase(repo)
}

func NewUserHandler(uc *usecase.UsuarioUseCase) *handler.UsuarioHandler {
	return handler.NewUserHandler(uc)
}
