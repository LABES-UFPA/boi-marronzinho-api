package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoacaoUseCase struct {
	doacaoRepo     repository.Repository[domain.Doacoes]
	itemDoacaoRepo repository.Repository[domain.ItemDoacao]
	userRepo       repository.UserRepository
	boicoinRepo    repository.BoicoinRepository
}

func NewDoacaoUseCase(
	doacaoRepo repository.Repository[domain.Doacoes],
	itemDoacaoRepo repository.Repository[domain.ItemDoacao],
	userRepo repository.UserRepository,
	boicoinRepo repository.BoicoinRepository,
) *DoacaoUseCase {
	return &DoacaoUseCase{
		doacaoRepo:     doacaoRepo,
		itemDoacaoRepo: itemDoacaoRepo,
		userRepo:       userRepo,
		boicoinRepo:    boicoinRepo,
	}
}

func (duc *DoacaoUseCase) AdicionaDoacao(doacaoRequest *domain.Doacoes) (*domain.Doacoes, error) {
	var resultDoacao *domain.Doacoes

	err := duc.doacaoRepo.WithTransaction(func(txRepo repository.Repository[domain.Doacoes]) error {
		itemDoacao, err := duc.itemDoacaoRepo.GetByID(doacaoRequest.ItemDoacaoID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("item de doação não encontrado")
			}
			return err
		}

		boicoinsRecebidos, err := calculaBoicoins(itemDoacao.BoicoinsPorUnidade, doacaoRequest.Quantidade)
		if err != nil {
			return err
		}

		doacao := &domain.Doacoes{
			ID:                uuid.New(),
			UsuarioID:         doacaoRequest.UsuarioID,
			ItemDoacaoID:      doacaoRequest.ItemDoacaoID,
			Quantidade:        doacaoRequest.Quantidade,
			BoicoinsRecebidos: boicoinsRecebidos,
			DataDoacao:        time.Now(),
		}

		createdDoacao, err := duc.doacaoRepo.Create(doacao)
		if err != nil {
			return err
		}

		// transacao := &domain.BoicoinsTransacoes{
		// 	ID:            uuid.New(),
		// 	UsuarioID:     doacaoRequest.UsuarioID,
		// 	Quantidade:    boicoinsRecebidos,
		// 	TipoTransacao: "recebimento_doacao",
		// 	Descricao:     "Recebimento de Boicoins por doação de item",
		// 	DataTransacao: time.Now(),
		// 	DoacaoID:      uuid.NullUUID{UUID: createdDoacao.ID, Valid: true},
		// 	PedidoID:      uuid.NullUUID{Valid: false},
		// 	PontoMapaID:   uuid.NullUUID{Valid: false},
		// }

		// if _, err := duc.boicoinRepo.Create(transacao); err != nil {
		// 	return err
		// }

		resultDoacao = createdDoacao
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultDoacao, nil
}

func (duc *DoacaoUseCase) CriarItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	itemDoacao := &domain.ItemDoacao{
		ID:                 uuid.New(),
		Descricao:          itemDoacaoRequest.Descricao,
		UnidadeMedida:      itemDoacaoRequest.UnidadeMedida,
		BoicoinsPorUnidade: itemDoacaoRequest.BoicoinsPorUnidade,
	}
	
	return duc.itemDoacaoRepo.Create(itemDoacao)
}

func (duc *DoacaoUseCase) AtualizaItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	itemDoacao, err := duc.itemDoacaoRepo.GetByID(itemDoacaoRequest.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item de doação não encontrado")
		}
		return nil, err
	}

	if itemDoacaoRequest.Descricao != "" {
		itemDoacao.Descricao = itemDoacaoRequest.Descricao
	}
	if itemDoacaoRequest.UnidadeMedida != "" {
		itemDoacao.UnidadeMedida = itemDoacaoRequest.UnidadeMedida
	}
	if itemDoacaoRequest.BoicoinsPorUnidade != 0 {
		itemDoacao.BoicoinsPorUnidade = itemDoacaoRequest.BoicoinsPorUnidade
	}

	return duc.itemDoacaoRepo.Update(itemDoacao)
}

func (duc *DoacaoUseCase) DeletarItemDoacao(id uuid.UUID) error {
	if err := duc.itemDoacaoRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item de doação não encontrado")
		}
		return err
	}
	return nil
}

func calculaBoicoins(valorUnidade float64, quantidade int64) (float64, error) {
	if valorUnidade < 0 || quantidade < 0 {
		return 0, errors.New("valorUnidade e quantidade devem ser maiores que zero")
	}
	return valorUnidade * float64(quantidade), nil
}
