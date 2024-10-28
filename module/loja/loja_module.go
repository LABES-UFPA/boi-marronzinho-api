package loja

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/core/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var LojaModule = fx.Options(
	fx.Provide(NewProdutoRepository),
	fx.Provide(NewPedidoRepository),
	fx.Provide(NewPedidoItensRepository),
	fx.Provide(NewCarrinhoRepository),
	fx.Provide(NewLojaUseCase),
	fx.Provide(NewLojaHandler),
)

func NewProdutoRepository(db *gorm.DB) repository.ProdutoRepository {
	return repository.NewProdutoRepository(db)
}

func NewPedidoRepository(db *gorm.DB) repository.PedidoRepository {
	return repository.NewPedidoRepository(db)
}

func NewPedidoItensRepository(db *gorm.DB) repository.PedidoItensRepository {
	return repository.NewPedidoItensRepository(db)
}

func NewCarrinhoRepository(db *gorm.DB) repository.CarrinhoRepository {
	return repository.NewCarrinhoRepository(db)
}

func NewLojaUseCase(
	produtoRepo repository.ProdutoRepository,
	pedidoRepo repository.PedidoRepository,
	pedidoItensRepo repository.PedidoItensRepository,
	carrinhoRepo repository.CarrinhoRepository,
	boicoinRepo repository.BoicoinRepository, // Já injetado a partir do módulo `boicoin`
) *usecase.LojaUseCase {
	return usecase.NewLojaUseCase(produtoRepo, pedidoRepo, pedidoItensRepo, carrinhoRepo, boicoinRepo)
}

func NewLojaHandler(luc *usecase.LojaUseCase) *handler.LojaHandler {
	return handler.NewLojaHandler(luc)
}
