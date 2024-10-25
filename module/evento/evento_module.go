package evento

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var EventoModule = fx.Options(
	fx.Provide(NewEventoRepository),
	fx.Provide(NewEventoUseCase),
	fx.Provide(NewEventoHandler),
)

func NewEventoRepository(db *gorm.DB) repository.EventoRepository {
	return repository.NewEventoRepository(db)
}

func NewEventoUseCase(eventoRepo repository.EventoRepository) *usecase.EventoUseCase {
	return usecase.NewEventoUseCase(eventoRepo)
}

func NewEventoHandler(uc *usecase.EventoUseCase) *handler.EventoHandler {
	return handler.NewEventoHandler(uc)
}
