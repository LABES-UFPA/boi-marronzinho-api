package repository

import (
    "boi-marronzinho-api/domain"
    "gorm.io/gorm"
)

type BoicoinRepository interface {
    Repository[domain.BoicoinsTransacoes]
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
