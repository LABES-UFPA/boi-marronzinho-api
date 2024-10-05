package domain

import "github.com/google/uuid"

type BoicoinsTransacoes struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID `json:"usuarioID"`
	Quantidade    float64   `json:"quantidade"`
	TipoTransacao string    `json:"tipoTransacao"`
	Descricao     string    `json:"descricao"`
	DataTransacao string    `json:"dataTransacao"`
	PedidoID      uuid.UUID `json:"pedidoId"`
	DoacaoID      uuid.UUID `json:"doacaoId"`
	PontoMapaID   uuid.UUID `json:"pontoMapaId"`
}

func (b *BoicoinsTransacoes) TableName() string {
	return "boi_marronzinho.boicoins_transacoes"
}
