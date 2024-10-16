package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// if err := oficinaRequest.Validate(); err != nil {
	// 	return nil, err
	// }

	oficina, err := o.oficinaRepo.Create(oficinaRequest)
	if err != nil {
		return nil, err
	}

	return oficina, nil
}

func (o *OficinaUseCase) ListaOficinas() ([]*domain.Oficinas, error) {
	return o.oficinaRepo.GetAll()
}

func (o *OficinaUseCase) UpdateOficina(id uuid.UUID, updateData *domain.Oficinas) (*domain.Oficinas, error) {
	oficina, err := o.oficinaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if updateData.Nome != "" {
		oficina.Nome = updateData.Nome
	}
	if updateData.LimiteParticipantes != nil {
		oficina.LimiteParticipantes = updateData.LimiteParticipantes
	}

	if _, err = o.oficinaRepo.Update(oficina); err != nil {
		return nil, errors.New("falha ao atualizar o usuário")
	}

	return oficina, nil
}

func (o *OficinaUseCase) DeleteUser(id uuid.UUID) error {
	_, err := o.oficinaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuário não encontrado")
		}
		return err
	}

	if err := o.oficinaRepo.Delete(id); err != nil {
		return errors.New("falha ao deletar oficia")
	}

	return nil
}

func (o *OficinaUseCase) InscricaoOficina(inscricaoRequest *domain.ParticipanteOficina, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error) {
	oficina, err := o.oficinaRepo.GetByID(inscricaoRequest.OficinaID)
	if err != nil {
		return nil, errors.New("oficina não encontrada")
	}

	if oficina.ParticipantesAtual >= *oficina.LimiteParticipantes {
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

func (o *OficinaUseCase) ListarTicketsPorUsuario(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error) {
	tickets, err := o.oficinaRepo.GetTicketsByUsuarioID(usuarioID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (o *OficinaUseCase) ValidaVoucher(codigoVoucher *string) (*dto.VoucherValidation, error) {

	validado, err := o.oficinaRepo.ValidaVoucher(codigoVoucher)
	if err != nil {
		return nil, err
	}

	validadoResponse := &dto.VoucherValidation{
		IDVoucher: *codigoVoucher,
		Validado:  validado.Validado,
	}

	return validadoResponse, nil
}
