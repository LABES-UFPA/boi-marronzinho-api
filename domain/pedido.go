package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Pedidos struct {
	ID             uuid.UUID  `json:"id" gorm:"primaryKey"`
	UsuarioID      uuid.UUID  `json:"usuarioId" validate:"required"`
	ProdutoID      uuid.UUID  `json:"produtoId" validate:"required"`
	OficinaID      uuid.UUID  `json:"oficinaId" validate:"required"`
	StatusPedido   string     `json:"statusPedido" validate:"required"` // Enum deve ser tratado em outro lugar
	EnderecoID     uuid.UUID  `json:"enderecoId" validate:"required"`
	PontoMapaID    uuid.UUID  `json:"pontoMapaId" validate:"required"`
	BoicoinsUsados float64    `json:"boicoinsUsados" gorm:"default:0.00" validate:"gte=0"`
	PrecoRealUsado float64    `json:"precoRealUsado" validate:"required,gte=0"`
	Quantidade     int        `json:"quantidade" gorm:"default:1" validate:"gte=1"`
	Codigo         string     `json:"codigo" validate:"required"`
	DataPedido     time.Time  `json:"dataPedido" gorm:"autoCreateTime"`
	DataConclusao  *time.Time `json:"dataConclusao"`

func (p *Pedidos) TableName() string {
	return "boi_marronzinho.pedidos"
}

func (p *Pedidos) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
