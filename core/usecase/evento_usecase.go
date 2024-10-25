package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventoUseCase struct {
	eventoRepo repository.EventoRepository
}

func NewEventoUseCase(eventoRepo repository.EventoRepository) *EventoUseCase {
	return &EventoUseCase{
		eventoRepo: eventoRepo,
	}
}

func (e *EventoUseCase) CriaEvento(eventoRequest *domain.Evento) (*domain.Evento, error) {
	evento, err := e.eventoRepo.Create(eventoRequest)
	if err != nil {
		return nil, err
	}

	

	return evento, nil
}

func (e *EventoUseCase) ListaEventos() ([]*domain.Evento, error) {
	return e.eventoRepo.GetAll()
}

func (e *EventoUseCase) GetEvento(id uuid.UUID) (*domain.Evento, error) {
	return e.eventoRepo.GetByID(id)
}

func (e *EventoUseCase) UpdateEvento(id uuid.UUID, updateData *domain.Evento) (*domain.Evento, error) {
	evento, err := e.eventoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("evento não encontrado")
		}
		return nil, err
	}

	if updateData.Nome != "" {
		evento.Nome = updateData.Nome
	}
	if updateData.Descricao != "" {
		evento.Descricao = updateData.Descricao
	}
	if updateData.LinkEndereco != "" {
		evento.LinkEndereco = updateData.LinkEndereco
	}
	if !updateData.DataEvento.IsZero() {
		evento.DataEvento = updateData.DataEvento
	}

	if _, err = e.eventoRepo.Update(evento); err != nil {
		return nil, errors.New("falha ao atualizar o evento")
	}

	return evento, nil
}

func (e *EventoUseCase) DeleteEvento(id uuid.UUID) error {
	_, err := e.eventoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("evento não encontrado")
		}
		return err
	}

	if err := e.eventoRepo.Delete(id); err != nil {
		return errors.New("falha ao deletar evento")
	}

	return nil
}
