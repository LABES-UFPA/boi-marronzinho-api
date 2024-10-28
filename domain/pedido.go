package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Pedidos struct {
	ID             uuid.UUID     `json:"id" gorm:"primaryKey"`
	UsuarioID      uuid.UUID     `json:"usuarioId" validate:"required"`
	StatusPedido   string        `json:"statusPedido" validate:"required"` // Enum deve ser tratado em outro lugar
	BoicoinsUsados float64       `json:"boicoinsUsados" gorm:"default:0.00" validate:"gte=0"`
	PrecoRealUsado float64       `json:"precoRealUsado" validate:"required,gte=0"`
	Quantidade     int           `json:"quantidade" gorm:"default:1" validate:"gte=1"`
	Codigo         string        `json:"codigo" validate:"required"`
	PedidoItens    []PedidoItens `json:"pedidoItens" gorm:"foreignKey:PedidoID"`
	DataPedido     time.Time     `json:"dataPedido" gorm:"autoCreateTime"`
	DataConclusao  *time.Time    `json:"dataConclusao"`
}

func (p *Pedidos) TableName() string {
	return "boi_marronzinho.pedidos"
}

func (p *Pedidos) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type PedidoItens struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	PedidoID      uuid.UUID `json:"pedidoId" gorm:"not null"`                // Referência ao pedido
	ProdutoID     uuid.UUID `json:"produtoId" gorm:"not null"`               // Referência ao produto
	Quantidade    int       `json:"quantidade" validate:"required,gte=1"`    // Quantidade de produtos
	PrecoUnitario float64   `json:"precoUnitario" validate:"required,gte=0"` // Preço unitário do produto
}

func (pi *PedidoItens) TableName() string {
	return "boi_marronzinho.pedido_itens"
}
