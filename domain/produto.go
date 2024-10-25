package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Produto struct {
	ID                  uuid.UUID `json:"id" gorm:"primaryKey"`
	Nome                string    `json:"nome" validate:"required"`
	Descricao           string    `json:"descricao" validate:"required"`
	PrecoBoicoins       float64   `json:"precoBoicoins" validate:"required"`
	PrecoReal           float64   `json:"precoReal" validate:"required"`
	QuantidadeEmEstoque int       `json:"quantidadeEmEstoque" validate:"required"`
	ImagemURL           string    `json:"imagemUrl"`
	CriadoEm            time.Time `json:"criadoEm" gorm:"autoCreateTime"`
}

func (p *Produto) TableName() string {
	return "boi_marronzinho.produtos"
}

func (p *Produto) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
