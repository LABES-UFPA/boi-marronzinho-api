package oficina

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var OficnaModule = fx.Options(
	fx.Provide(NewOficinaRepository),
	fx.Provide(NewOficinaUseCase),
	fx.Provide(NewOficinaHandler),
)

func NewOficinaRepository(db *gorm.DB) repository.OficinaRepository {
	return repository.NewOficinaRepository(db)
}

func NewOficinaUseCase(oficinaRepo repository.OficinaRepository, usuarioRepo repository.UserRepository) *usecase.OficinaUseCase {
	return usecase.NewOficinaUseCase(oficinaRepo, usuarioRepo)
}

func NewOficinaHandler(uc *usecase.OficinaUseCase) *handler.OficinaHandler {
	return handler.NewOficinaHandler(uc)
}
