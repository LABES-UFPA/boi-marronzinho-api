package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Buscando saldo de Boicoins para o usuário com ID: %s", usuarioID)
	var saldoTotal float32
	err := r.db.Table("boi_marronzinho.boicoins_transacoes").
		Select("SUM(quantidade)").
		Where("usuario_id = ?", usuarioID).
		Scan(&saldoTotal).Error

	if err != nil {
		logrus.Errorf("Erro ao buscar saldo de Boicoins para o usuário com ID %s: %v", usuarioID, err)
		return 0, err
	}

	logrus.Infof("Saldo de Boicoins encontrado para o usuário com ID %s: %f", usuarioID, saldoTotal)
	return saldoTotal, nil
}
