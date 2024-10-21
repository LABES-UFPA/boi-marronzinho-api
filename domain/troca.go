package domain

import (
	"time"

	"github.com/google/uuid"
)

type Troca struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID         uuid.UUID `json:"usuarioId"`
	ItemTrocaID       uuid.UUID `json:"itemTrocaId"`
	Quantidade        float64   `json:"quantidade"`
	BoicoinsRecebidos float64   `json:"boicoinsRecebidos"`
	Status            string    `json:"status"` // Novo campo para status ("pendente", "validada", "rejeitada")
	DataDoacao        time.Time `json:"dataDoacao"`
}

func (t *Troca) TableName() string {
	return "boi_marronzinho.troca"
}

type ItemTroca struct {
	ID                 uuid.UUID `json:"id" gorm:"primaryKey"`
	NomeItem           string    `json:"nomeItem"`
	Descricao          string    `json:"Descricao"`
	UnidadeMedida      string    `json:"unidadeMedida"`
	BoicoinsPorUnidade float64   `json:"boicoinsPorUnidade"`
}

func (it *ItemTroca) TableName() string {
	return "boi_marronzinho.item_troca"
}
