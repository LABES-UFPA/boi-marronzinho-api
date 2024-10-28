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
	produtoRepo     repository.ProdutoRepository
	pedidoItensRepo repository.PedidoItensRepository
	pedidoRepo      repository.PedidoRepository
	carrinhoRepo    repository.CarrinhoRepository
	transacaoRepo   repository.BoicoinRepository
}

func NewLojaUseCase(produtoRepo repository.ProdutoRepository, pedidoRepo repository.PedidoRepository, pedidoItensRepo repository.PedidoItensRepository, carrinhoRepo repository.CarrinhoRepository, transacaoRepo repository.BoicoinRepository) *LojaUseCase {
	return &LojaUseCase{
		produtoRepo:     produtoRepo,
		pedidoRepo:      pedidoRepo,
		pedidoItensRepo: pedidoItensRepo,
		carrinhoRepo:    carrinhoRepo,
		transacaoRepo:   transacaoRepo,
	}
}

func (l *LojaUseCase) AdicionarOuIncrementarItemCarrinho(usuarioID, produtoID uuid.UUID, quantidade int) error {
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

	// Verifica se o item já existe no carrinho
	itemExistente, err := l.carrinhoRepo.GetByUsuarioEProdutoID(usuarioID, produtoID)
	if err != nil {
		return err
	}

	if itemExistente != nil {
		// Incrementa a quantidade do item existente no carrinho
		itemExistente.Quantidade += quantidade
		if _, err := l.carrinhoRepo.Update(itemExistente); err != nil {
			return err
		}
	} else {
		// Adiciona novo item ao carrinho
		itemCarrinho := &domain.CarrinhoItem{
			ID:            uuid.New(),
			UsuarioID:     usuarioID,
			ProdutoID:     produtoID,
			Quantidade:    quantidade,
			PrecoUnitario: produto.PrecoBoicoins,
		}
		if _, err := l.carrinhoRepo.Create(itemCarrinho); err != nil {
			return err
		}
	}
	return nil
}

func (l *LojaUseCase) ListarItensCarrinho(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error) {
	return l.carrinhoRepo.GetByUsuarioID(usuarioID)
}

func (l *LojaUseCase) RemoverItemCarrinho(usuarioID, itemID uuid.UUID) error {
	return l.carrinhoRepo.Delete(itemID)
}

func (l *LojaUseCase) AtualizarQuantidadeItemCarrinho(usuarioID, produtoID uuid.UUID, quantidade int) error {
	itemExistente, err := l.carrinhoRepo.GetByUsuarioEProdutoID(usuarioID, produtoID)
	if err != nil {
		return err
	}

	if itemExistente == nil {
		return errors.New("item não encontrado no carrinho")
	}

	if quantidade <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}

	// Verifica o estoque do produto antes de atualizar
	produto, err := l.produtoRepo.GetByID(produtoID)
	if err != nil {
		return err
	}

	if produto.QuantidadeEmEstoque < quantidade {
		return errors.New("quantidade solicitada excede o estoque disponível")
	}

	itemExistente.Quantidade = quantidade
	_, err = l.carrinhoRepo.Update(itemExistente)
	return err
}

func (l *LojaUseCase) FinalizarCompra(usuarioID uuid.UUID) (*domain.Pedidos, error) {
	itensCarrinho, err := l.carrinhoRepo.GetByUsuarioID(usuarioID)
	if err != nil || len(itensCarrinho) == 0 {
		return nil, errors.New("carrinho vazio ou erro ao buscar itens")
	}

	produtoIDs := make([]uuid.UUID, len(itensCarrinho))
	for i, item := range itensCarrinho {
		produtoIDs[i] = item.ProdutoID
	}

	produtos, err := l.produtoRepo.GetByIDs(produtoIDs)
	if err != nil {
		return nil, errors.New("erro ao buscar dados dos produtos")
	}

	produtoMap := make(map[uuid.UUID]*domain.Produto)
	for _, produto := range produtos {
		produtoMap[produto.ID] = produto
	}

	var totalBoicoins float64
	for _, item := range itensCarrinho {
		produto, ok := produtoMap[item.ProdutoID]
		if !ok || produto.QuantidadeEmEstoque < item.Quantidade {
			return nil, errors.New("estoque insuficiente para um ou mais itens no carrinho")
		}
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
		StatusPedido:   "concluido",
		DataPedido:     time.Now(),
	}

	createdPedido, err := l.pedidoRepo.Create(pedido)
	if err != nil {
		return nil, errors.New("erro ao criar pedido")
	}

	transacao := &domain.BoicoinsTransacoes{
		ID:            uuid.New(),
		UsuarioID:     usuarioID,
		Quantidade:    -totalBoicoins,
		TipoTransacao: "compra_produto",
		DataTransacao: time.Now(),
		PedidoID:      &pedido.ID,
	}

	createdTransacao, err := l.transacaoRepo.Create(transacao)
	if err != nil {
		_ = l.pedidoRepo.Delete(createdPedido.ID)
		return nil, errors.New("erro ao registrar transação de Boicoins")
	}

	var pedidoItens []*domain.PedidoItens
	for _, item := range itensCarrinho {
		pedidoItem := &domain.PedidoItens{
			ID:            uuid.New(),
			PedidoID:      createdPedido.ID,
			ProdutoID:     item.ProdutoID,
			Quantidade:    item.Quantidade,
			PrecoUnitario: item.PrecoUnitario,
		}
		pedidoItens = append(pedidoItens, pedidoItem)

		produto := produtoMap[item.ProdutoID]
		produto.QuantidadeEmEstoque -= item.Quantidade
	}

	for _, item := range pedidoItens {
		if _, err := l.pedidoItensRepo.Create(item); err != nil {
			_ = l.pedidoRepo.Delete(createdPedido.ID)
			_ = l.transacaoRepo.Delete(createdTransacao.ID)
			return nil, errors.New("erro ao criar itens do pedido")
		}
	}

	if err := l.produtoRepo.BatchUpdate(produtos); err != nil {
		_ = l.pedidoRepo.Delete(createdPedido.ID)
		_ = l.transacaoRepo.Delete(createdTransacao.ID)
		for _, item := range pedidoItens {
			_ = l.pedidoRepo.Delete(item.ID)
		}
		return nil, errors.New("erro ao atualizar estoque dos produtos")
	}

	if err := l.carrinhoRepo.BatchDeleteByIDs(itensCarrinho); err != nil {
		_ = l.pedidoRepo.Delete(createdPedido.ID)
		_ = l.transacaoRepo.Delete(createdTransacao.ID)
		for _, item := range pedidoItens {
			_ = l.pedidoRepo.Delete(item.ID)
		}
		return nil, errors.New("erro ao limpar carrinho")
	}

	return createdPedido, nil
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
			return errors.New("produto não encontrado")
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
