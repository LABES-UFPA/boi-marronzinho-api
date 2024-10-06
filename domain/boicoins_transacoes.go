package domain

import (
	"time"

	"github.com/google/uuid"
)

type BoicoinsTransacoes struct {
	ID            uuid.UUID     `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID     `json:"usuarioID"`
	Quantidade    float64       `json:"quantidade"`
	TipoTransacao string        `json:"tipoTransacao"`
	Descricao     string        `json:"descricao"`
	DataTransacao time.Time     `json:"dataTransacao"`
	PedidoID      uuid.NullUUID `json:"pedidoId"`
	DoacaoID      uuid.NullUUID `json:"doacaoId"`
	PontoMapaID   uuid.NullUUID `json:"pontoMapaId"`
}

func (b *BoicoinsTransacoes) TableName() string {
	return "boi_marronzinho.boicoins_transacoes"
}
