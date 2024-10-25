package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CarrinhoItem struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID `json:"usuarioId" validate:"required"`
	ProdutoID     uuid.UUID `json:"produtoId" validate:"required"`
	Quantidade    int       `json:"quantidade" validate:"required,gte=1"`
	PrecoUnitario float64   `json:"precoUnitario" validate:"required,gte=0"`
	CriadoEm      time.Time `json:"criadoEm" gorm:"autoCreateTime"`
}

func (c *CarrinhoItem) TableName() string {
	return "boi_marronzinho.carrinho_itens"
}

func (c *CarrinhoItem) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
