package doacao

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var DoacaoModule = fx.Options(
	fx.Provide(NewDoacaoRepository),
	fx.Provide(NewDoacaoUseCase),
	fx.Provide(NewDoacaoHandler),
)

func NewDoacaoRepository(db *gorm.DB) repository.DoacaoRepository {
	return repository.NewDoacaoRepository(db)
}

func NewDoacaoUseCase(repo repository.DoacaoRepository) *usecase.DoacaoUseCase {
	return usecase.NewDoacaoUseCase(repo)
}

func NewDoacaoHandler(duc *usecase.DoacaoUseCase) *handler.DoacaoHandler {
	return handler.NewDoacaoHandler(duc)
}
