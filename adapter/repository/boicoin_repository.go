package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoicoinRepository interface {
	Repository[domain.BoicoinsTransacoes]
	GetSaldoBoicoins(usuarioID uuid.UUID) (float32, error)
}

type boicoinRepository struct {
	Repository[domain.BoicoinsTransacoes]
	db *gorm.DB
}

func NewBoicoinRepository(db *gorm.DB) BoicoinRepository {
	return &boicoinRepository{
		Repository: NewRepository[domain.BoicoinsTransacoes](db),
		db:         db,
	}
}

func (r *boicoinRepository) GetSaldoBoicoins(usuarioID uuid.UUID) (float32, error) {
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