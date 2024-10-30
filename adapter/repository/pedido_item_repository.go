package repository

import (
	"boi-marronzinho-api/domain"

	"gorm.io/gorm"
)

type PedidoItensRepository interface {
	Repository[domain.PedidoItens]
}

type pedidoItensRepository struct {
	Repository[domain.PedidoItens]
	db *gorm.DB
}

func NewPedidoItensRepository(db *gorm.DB) PedidoItensRepository {
	return &pedidoItensRepository{
		Repository: NewRepository[domain.PedidoItens](db),
		db:         db,
	}
}
