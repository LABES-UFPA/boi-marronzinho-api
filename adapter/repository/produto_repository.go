package repository

import (
	"boi-marronzinho-api/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProdutoRepository interface {
	Repository[domain.Produto]
	GetByIDs(ids []uuid.UUID) ([]*domain.Produto, error)
}

type produtoRepository struct {
	Repository[domain.Produto]
	db *gorm.DB
}

func NewProdutoRepository(db *gorm.DB) ProdutoRepository {
	return &produtoRepository{
		Repository: NewRepository[domain.Produto](db),
		db:         db,
	}
}

func (r *produtoRepository) GetByIDs(ids []uuid.UUID) ([]*domain.Produto, error) {
	logrus.Infof("Buscando múltiplos produtos por IDs")
	var produtos []*domain.Produto
	if err := r.db.Where("id IN ?", ids).Find(&produtos).Error; err != nil {
		logrus.Errorf("Erro ao buscar múltiplos produtos: %v", err)
		return nil, err
	}
	logrus.Infof("Encontrados %d produtos", len(produtos))
	return produtos, nil
}
