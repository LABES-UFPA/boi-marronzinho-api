package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CarrinhoRepository interface {
	Repository[domain.CarrinhoItem]
	GetByUsuarioID(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error)
	GetSaldoBoicoins(usuarioID uuid.UUID) (float32, error)
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
	logrus.Infof("Buscando itens do carrinho para o usuário com ID: %s", usuarioID)
	var itens []*domain.CarrinhoItem
	err := r.db.Where("usuario_id = ?", usuarioID).Find(&itens).Error
	if err != nil {
		logrus.Errorf("Erro ao buscar itens do carrinho para o usuário com ID %s: %v", usuarioID, err)
		return nil, err
	}
	logrus.Infof("Encontrados %d itens no carrinho para o usuário com ID: %s", len(itens), usuarioID)
	return itens, nil
}

func (r *carrinhoRepository) GetSaldoBoicoins(usuarioID uuid.UUID) (float32, error) {
	logrus.Infof("Buscando saldo de Boicoins para o usuário com ID: %s", usuarioID)
	var saldoTotal float32
	err := r.db.Table("boicoins_transacoes").
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
