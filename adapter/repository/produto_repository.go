package repository

import (
	"boi-marronzinho-api/domain"

	"gorm.io/gorm"
)

type ProdutoRepository interface {
	Repository[domain.Produto]
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
