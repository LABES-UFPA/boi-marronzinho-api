package boicoin

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var BoicoinModule = fx.Options(
    fx.Provide(NewBoicoinRepository),
    fx.Provide(NewCarteiraHandler),
    fx.Provide(NewCarteiraUseCase),
)

func NewBoicoinRepository(db *gorm.DB) repository.BoicoinRepository {
    return repository.NewBoicoinRepository(db)
}

func NewCarteiraHandler(cuc *usecase.CarteiraUseCase) *handler.CarteiraHandler {
    return handler.NewCarteiraHandler(cuc)
}

func NewCarteiraUseCase(
    boicoinRepo repository.BoicoinRepository,
) *usecase.CarteiraUseCase {
    return usecase.NewCarteiraUseCase(boicoinRepo)
}
