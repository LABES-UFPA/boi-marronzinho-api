package repository

import (
	"boi-marronzinho-api/domain"

	"gorm.io/gorm"
)

type EventoRepository interface {
	Repository[domain.Evento]
}

type eventoRepository struct {
	Repository[domain.Evento]
	db *gorm.DB
}

func NewEventoRepository(db *gorm.DB) EventoRepository {
	return &eventoRepository{
		Repository: NewRepository[domain.Evento](db),
		db:         db,
	}
}
