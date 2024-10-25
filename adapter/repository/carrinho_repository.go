package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CarrinhoRepository interface {
	Repository[domain.CarrinhoItem]
	GetByUsuarioID(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error)
}

type carrinhoRepository struct {
	Repository[domain.CarrinhoItem]
	db *gorm.DB
}

func NewCarrinhoRepository(db *gorm.DB) CarrinhoRepository {
	return &carrinhoRepository{
		Repository: NewRepository[domain.CarrinhoItem](db),
		db:         db,
	}
}

func (r *carrinhoRepository) GetByUsuarioID(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error) {
	var itens []*domain.CarrinhoItem
	err := r.db.Where("usuario_id = ?", usuarioID).Find(&itens).Error
	if err != nil {
		return nil, err
	}
	return itens, nil
}

func (r *carrinhoRepository) GetSaldoBoicoins(usuarioID uuid.UUID) (float32, error) {
	var saldoTotal float32
	err := r.db.Table("boicoins_transacoes").
		Select("SUM(quantidade)").
		Where("usuario_id = ?", usuarioID).
		Scan(&saldoTotal).Error

	if err != nil {
		return 0, err
	}

	return saldoTotal, nil
}
