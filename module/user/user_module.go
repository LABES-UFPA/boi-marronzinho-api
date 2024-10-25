package user

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

func NewUserUseCase(userRepo repository.UserRepository, boicoinRepo repository.BoicoinRepository) *usecase.UserUseCase {
	return usecase.NewUsuarioUseCase(userRepo, boicoinRepo)
}

func NewUserHandler(uc *usecase.UserUseCase) *handler.UserHandler {
	return handler.NewUserHandler(uc)
}
