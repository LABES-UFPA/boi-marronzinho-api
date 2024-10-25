package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LojaUseCase struct {
	produtoRepo   repository.ProdutoRepository
	pedidoRepo    repository.PedidoRepository
	carrinhoRepo  repository.CarrinhoRepository
	transacaoRepo repository.BoicoinRepository
}

func NewLojaUseCase(produtoRepo repository.ProdutoRepository, pedidoRepo repository.PedidoRepository, carrinhoRepo repository.CarrinhoRepository, transacaoRepo repository.BoicoinRepository) *LojaUseCase {
	return &LojaUseCase{
		produtoRepo:   produtoRepo,
		pedidoRepo:    pedidoRepo,
		carrinhoRepo:  carrinhoRepo,
		transacaoRepo: transacaoRepo,
	}
}

func (l *LojaUseCase) AdicionarItemCarrinho(usuarioID, produtoID uuid.UUID, quantidade int) error {
	produto, err := l.produtoRepo.GetByID(produtoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produto não encontrado")
		}
		return err
	}

	if produto.QuantidadeEmEstoque < quantidade {
		return errors.New("quantidade solicitada excede o estoque disponível")
	}

	itemCarrinho := &domain.CarrinhoItem{
		ID:            uuid.New(),
		UsuarioID:     usuarioID,
		ProdutoID:     produtoID,
		Quantidade:    quantidade,
		PrecoUnitario: produto.PrecoBoicoins,
	}

	_, err = l.carrinhoRepo.Create(itemCarrinho)
	if err != nil {
		return err
	}

	return nil
}

func (l *LojaUseCase) ListarItensCarrinho(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error) {
	return l.carrinhoRepo.GetByUsuarioID(usuarioID)
}

func (l *LojaUseCase) RemoverItemCarrinho(usuarioID, itemID uuid.UUID) error {
	return l.carrinhoRepo.Delete(itemID)
}

func (l *LojaUseCase) FinalizarCompra(usuarioID uuid.UUID) (*domain.Pedidos, error) {
	itensCarrinho, err := l.carrinhoRepo.GetByUsuarioID(usuarioID)
	if err != nil || len(itensCarrinho) == 0 {
		return nil, errors.New("carrinho vazio ou erro ao buscar itens")
	}

	var totalBoicoins float64
	for _, item := range itensCarrinho {
		totalBoicoins += item.PrecoUnitario * float64(item.Quantidade)
	}

	saldoAtual, err := l.transacaoRepo.GetSaldoBoicoins(usuarioID)
	if err != nil || saldoAtual < float32(totalBoicoins) {
		return nil, errors.New("saldo insuficiente em Boicoins")
	}

	pedido := &domain.Pedidos{
		ID:             uuid.New(),
		UsuarioID:      usuarioID,
		BoicoinsUsados: totalBoicoins,
		StatusPedido:   "concluído",
		DataPedido:     time.Now(),
	}

	_, err = l.pedidoRepo.Create(pedido)
	if err != nil {
		return nil, err
	}

	transacao := &domain.BoicoinsTransacoes{
		ID:            uuid.New(),
		UsuarioID:     usuarioID,
		Quantidade:    -totalBoicoins,
		TipoTransacao: "compra_produto",
		DataTransacao: time.Now(),
		PedidoID:      &pedido.ID,
	}

	_, err = l.transacaoRepo.Create(transacao)
	if err != nil {
		return nil, err
	}

	for _, item := range itensCarrinho {
		produto, _ := l.produtoRepo.GetByID(item.ProdutoID)
		produto.QuantidadeEmEstoque -= item.Quantidade
		l.produtoRepo.Update(produto)
		l.carrinhoRepo.Delete(item.ID)
	}

	return pedido, nil
}

func (l *LojaUseCase) AdicionarProduto(produtoRequest *domain.Produto) (*domain.Produto, error) {
	if err := produtoRequest.Validate(); err != nil {
		return nil, err
	}

	produtoRequest.ID = uuid.New()

	produto, err := l.produtoRepo.Create(produtoRequest)
	if err != nil {
		return nil, err
	}

	return produto, nil
}

func (l *LojaUseCase) ListaProdutos() ([]*domain.Produto, error) {
	return l.produtoRepo.GetAll()
}

func (l *LojaUseCase) RemoveProduto(id uuid.UUID) error {
	_, err := l.produtoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuário não encontrado")
		}
		return err
	}

	if err := l.produtoRepo.Delete(id); err != nil {
		return errors.New("falha ao deletar o produto")
	}

	return nil
}

func (l *LojaUseCase) EditaProduto(id uuid.UUID, updateData *domain.Produto) (*domain.Produto, error) {
	produto, err := l.produtoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}

	if updateData.Nome != "" {
		produto.Nome = updateData.Nome
	}
	if updateData.Descricao != "" {
		produto.Descricao = updateData.Descricao
	}
	if updateData.PrecoBoicoins > 0 {
		produto.PrecoBoicoins = updateData.PrecoBoicoins
	}
	if updateData.PrecoReal > 0 {
		produto.PrecoReal = updateData.PrecoReal
	}

	if _, err = l.produtoRepo.Update(produto); err != nil {
		return nil, errors.New("falha ao atualizar o produto")
	}

	return produto, nil
}
