package container

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/postgres"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Função para inicializar o banco de dados
func NewDB() *gorm.DB {
	return postgres.InitDB()
}

// Função para inicializar o UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

// Função para inicializar o UserUseCase
func NewUserUseCase(repo repository.UserRepository) *usecase.UsuarioUseCase {
	return usecase.NewUsuarioUseCase(repo)
}

// Função para inicializar o UserHandler
func NewUserHandler(uc *usecase.UsuarioUseCase) *handler.UsuarioHandler {
	return handler.NewUserHandler(uc)
}

// Define o container para uso no fx
func NewContainer() fx.Option {
	return fx.Options(
		fx.Provide(
			NewDB,             // Provedor para o banco de dados
			NewUserRepository, // Provedor para o repositório
			NewUserUseCase,    // Provedor para o caso de uso
			NewUserHandler,    // Provedor para o handler
		),
	)
}
