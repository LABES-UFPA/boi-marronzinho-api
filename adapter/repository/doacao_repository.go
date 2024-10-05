package repository

import (
    "boi-marronzinho-api/domain"
    "gorm.io/gorm"
)

type DoacaoRepository interface {
    Repository[domain.ItemDoacao]
}

type doacaoRepository struct {
    Repository[domain.ItemDoacao]
    db *gorm.DB
}

func NewDoacaoRepository(db *gorm.DB) DoacaoRepository {
    return &doacaoRepository{
        Repository: NewRepository[domain.ItemDoacao](db),
        db: db,
    }
}
