package doacao

import (
    "boi-marronzinho-api/adapter/http/handler"
    "boi-marronzinho-api/adapter/repository"
    "boi-marronzinho-api/core/usecase"
    "boi-marronzinho-api/domain"

    "go.uber.org/fx"
    "gorm.io/gorm"
)

var DoacaoModule = fx.Options(
    fx.Provide(NewDoacaoRepository),
    fx.Provide(NewItemDoacaoRepository),
    fx.Provide(NewDoacaoUseCase),
    fx.Provide(NewDoacaoHandler),
)

func NewDoacaoRepository(db *gorm.DB) repository.Repository[domain.Doacoes] {
    return repository.NewRepository[domain.Doacoes](db)
}

func NewItemDoacaoRepository(db *gorm.DB) repository.Repository[domain.ItemDoacao] {
    return repository.NewRepository[domain.ItemDoacao](db)
}

func NewDoacaoUseCase(doacaoRepo repository.Repository[domain.Doacoes], itemDoacaoRepo repository.Repository[domain.ItemDoacao]) *usecase.DoacaoUseCase {
    return usecase.NewDoacaoUseCase(doacaoRepo, itemDoacaoRepo)
}

func NewDoacaoHandler(duc *usecase.DoacaoUseCase) *handler.DoacaoHandler {
    return handler.NewDoacaoHandler(duc)
}
