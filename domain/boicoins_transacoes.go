package domain

import (
	"time"

	"github.com/google/uuid"
)

type BoicoinsTransacoes struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID  `json:"usuarioID"`
	Quantidade    float64    `json:"quantidade"`
	TipoTransacao string     `json:"tipoTransacao"`
	Descricao     string     `json:"descricao"`
	DataTransacao time.Time  `json:"dataTransacao"`
	PedidoID      *uuid.UUID `json:"pedidoId"`
	TrocaID       *uuid.UUID `json:"trocaId"`
	OficinaID     *uuid.UUID `json:"oficinaId"`
}

func (b *BoicoinsTransacoes) TableName() string {
	return "boi_marronzinho.boicoins_transacoes"
}
