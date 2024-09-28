package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
)

type DoacaoUseCase struct {
	doacaoRepo repository.DoacaoRepository
}

func NewDoacaoUseCase(doacaoRepo repository.DoacaoRepository) *DoacaoUseCase {
	return &DoacaoUseCase{doacaoRepo: doacaoRepo}
}

func (duc *DoacaoUseCase) AdicionaDoacao(doacaoResquest *domain.Doacoes) (*domain.Doacoes, error) {

	itemDoacao, err := duc.doacaoRepo.CapturaItemDoacao(doacaoResquest.ItemDoacaoID)
	if err != nil {
		return nil, nil
	}

	boicoinsRecebidos, err := calculaBoicoins(itemDoacao.BoicoinsUnidade, doacaoResquest.Quantidade)
	if err != nil {
		return nil, nil
	}

	doacao := &domain.Doacoes{
		ID:                uuid.New(),
		UsuarioID:         doacaoResquest.UsuarioID,
		ItemDoacaoID:      doacaoResquest.ItemDoacaoID,
		BoicoinsRecebidos: boicoinsRecebidos,
		DataDoacao:        time.Now(),
	}
	return doacao, nil
}

func calculaBoicoins(valorUnidade float64, quantidade int64) (float64, error) {
	if valorUnidade < 0 || quantidade < 0 {
		return 0, errors.New("valorUnidade e quantidade devem ser maiores que zero")
	}
	return valorUnidade * float64(quantidade), nil
}
