package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
)

type CarteiraUseCase struct {
	boicoinRepo repository.BoicoinRepository
}

func NewCarteiraUseCase(boicoinRepo repository.BoicoinRepository) *CarteiraUseCase {
	return &CarteiraUseCase{
		boicoinRepo: boicoinRepo,
	}
}

func (cuc *CarteiraUseCase) CriaTransacao(transacao *domain.BoicoinsTransacoes) error {
	_, err := cuc.boicoinRepo.Create(transacao)

	return err
}
