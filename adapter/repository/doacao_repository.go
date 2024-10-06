package repository

import (
    "boi-marronzinho-api/domain"
    "gorm.io/gorm"
)

type DoacaoRepository interface {
    Repository[domain.Doacoes]
}

type doacaoRepository struct {
    Repository[domain.Doacoes]
    db *gorm.DB
}

func NewDoacaoRepository(db *gorm.DB) DoacaoRepository {
    return &doacaoRepository{
        Repository: NewRepository[domain.Doacoes](db),
        db:         db,
    }
}
