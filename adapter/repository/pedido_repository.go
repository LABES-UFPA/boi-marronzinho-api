package repository

import (
	"boi-marronzinho-api/domain"

	"gorm.io/gorm"
)

type PedidoRepository interface {
	Repository[domain.Pedidos]
}

type pedidoRepository struct {
	Repository[domain.Pedidos]
	db *gorm.DB
}

func NewPedidoRepository(db *gorm.DB) PedidoRepository {
	return &pedidoRepository{
		Repository: NewRepository[domain.Pedidos](db),
		db:         db,
	}
}
