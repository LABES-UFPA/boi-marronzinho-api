package repository

import (
	"boi-marronzinho-api/domain"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CarrinhoRepository interface {
	Repository[domain.CarrinhoItem]
	GetByUsuarioID(usuarioID uuid.UUID) ([]*domain.CarrinhoItem, error)
	GetByUsuarioEProdutoID(usuarioID, produtoID uuid.UUID) (*domain.CarrinhoItem, error)
	BatchDeleteByIDs(items []*domain.CarrinhoItem) error
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
	if err := r.db.Where("usuario_id = ?", usuarioID).Find(&itens).Error; err != nil {
		logrus.Errorf("Erro ao buscar itens do carrinho para o usuário com ID %s: %v", usuarioID, err)
		return nil, err
	}
	logrus.Infof("Encontrados %d itens no carrinho para o usuário com ID: %s", len(itens), usuarioID)
	return itens, nil
}

func (r *carrinhoRepository) GetByUsuarioEProdutoID(usuarioID, produtoID uuid.UUID) (*domain.CarrinhoItem, error) {
	logrus.Infof("Buscando item do carrinho para o usuário com ID: %s e produto com ID: %s", usuarioID, produtoID)
	var item domain.CarrinhoItem
	if err := r.db.Where("usuario_id = ? AND produto_id = ?", usuarioID, produtoID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Infof("Item do carrinho não encontrado para o usuário com ID: %s e produto com ID: %s", usuarioID, produtoID)
			return nil, nil
		}
		logrus.Errorf("Erro ao buscar item do carrinho: %v", err)
		return nil, err
	}
	logrus.Infof("Item do carrinho encontrado: %+v", item)
	return &item, nil
}

func (r *carrinhoRepository) BatchDeleteByIDs(items []*domain.CarrinhoItem) error {
	logrus.Infof("Deletando múltiplos itens do carrinho")
	ids := make([]uuid.UUID, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	if err := r.db.Delete(&domain.CarrinhoItem{}, ids).Error; err != nil {
		logrus.Errorf("Erro ao deletar múltiplos itens do carrinho: %v", err)
		return err
	}
	logrus.Infof("Deleção em batch dos itens do carrinho concluída com sucesso")
	return nil
}
