package repository

import (
	"boi-marronzinho-api/domain"

	"gorm.io/gorm"
)

type TrocaRepository interface {
	Repository[domain.Troca]
}

type trocaRepository struct {
	Repository[domain.Troca]
	db *gorm.DB
}

func NewTocaRepository(db *gorm.DB) TrocaRepository {
	return &trocaRepository{
		Repository: NewRepository[domain.Troca](db),
		db:         db,
	}
}
