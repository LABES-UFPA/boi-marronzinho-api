package troca

import (
    "boi-marronzinho-api/adapter/http/handler"
    "boi-marronzinho-api/adapter/repository"
    "boi-marronzinho-api/core/usecase"
    "boi-marronzinho-api/domain"

    "go.uber.org/fx"
    "gorm.io/gorm"
)

var TrocaModule = fx.Options(
    fx.Provide(NewTrocaRepository),
    fx.Provide(NewItemTrocaRepository),
    fx.Provide(NewTrocaUseCase),
    fx.Provide(NewTrocaHandler),
)

func NewTrocaRepository(db *gorm.DB) repository.TrocaRepository {
    return repository.NewTrocaRepository(db)
}

func NewItemTrocaRepository(db *gorm.DB) repository.Repository[domain.ItemTroca] {
    return repository.NewRepository[domain.ItemTroca](db)
}

func NewTrocaUseCase(
    trocaRepo repository.TrocaRepository,
    ItemTrocaRepo repository.Repository[domain.ItemTroca],
    usuarioRepo repository.UserRepository,
    boicoinRepo repository.BoicoinRepository,
) *usecase.TrocaUseCase {
    return usecase.NewTrocaUseCase(trocaRepo, ItemTrocaRepo, usuarioRepo, boicoinRepo)
}

func NewTrocaHandler(tuc *usecase.TrocaUseCase) *handler.TrocaHandler {
    return handler.NewTrocaHandler(tuc)
}
