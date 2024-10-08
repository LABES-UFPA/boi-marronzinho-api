package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
)

type OficinaUseCase struct {
	oficinaRepo repository.OficinaRepository
	usuarioRepo repository.UserRepository
}

func NewOficinaUseCase(oficinaRepo repository.OficinaRepository, usuarioRepo repository.UserRepository) *OficinaUseCase {
	return &OficinaUseCase{
		oficinaRepo: oficinaRepo,
		usuarioRepo: usuarioRepo,
	}
}

func (o *OficinaUseCase) CriaOficina(oficinaRequest *domain.Oficinas) (*domain.Oficinas, error) {
	if err := oficinaRequest.Validate(); err != nil {
		return nil, err
	}

	oficina, err := o.oficinaRepo.Create(oficinaRequest)
	if err != nil {
		return nil, err
	}

	return oficina, nil
}

func (o *OficinaUseCase) ListaOficinas() ([]*domain.Oficinas, error) {
	return o.oficinaRepo.GetAll()
}
func (o *OficinaUseCase) InscricaoOficina(inscricaoRequest *domain.ParticipanteOficina, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error) {
	oficina, err := o.oficinaRepo.GetByID(inscricaoRequest.OficinaID)
	if err != nil {
		return nil, errors.New("oficina não encontrada")
	}

	if oficina.ParticipantesAtual >= oficina.LimiteParticipantes {
		return nil, errors.New("não há mais vagas disponíveis para esta oficina")
	}

	usuario, err := o.usuarioRepo.GetByID(inscricaoRequest.UsuarioID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if pagamentoEmBoicoins {
		if usuario.SaldoBoicoins < float32(oficina.PrecoBoicoins) {
			return nil, errors.New("saldo de Boicoins insuficiente")
		}
	}

	err = o.oficinaRepo.InscreverParticipante(inscricaoRequest.OficinaID, usuario, pagamentoEmBoicoins)
	if err != nil {
		return nil, err
	}

	return inscricaoRequest, nil
}
