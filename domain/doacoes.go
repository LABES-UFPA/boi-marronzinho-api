package domain

import (
	"time"

	"github.com/google/uuid"
)

type Doacoes struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID         uuid.UUID `json:"usuarioId"`
	ItemDoacaoID      uuid.UUID `json:"itemDoacaoId"`
	Quantidade        float64   `json:"quantidade"`
	BoicoinsRecebidos float64   `json:"boicoinsRecebidos"`
	Status            string    `json:"status"` // Novo campo para status ("pendente", "validada", "rejeitada")
	DataDoacao        time.Time `json:"dataDoacao"`
}

func (d *Doacoes) TableName() string {
	return "boi_marronzinho.doacoes"
}

type ItemDoacao struct {
	ID                 uuid.UUID `json:"id" gorm:"primaryKey"`
	NomeItem           string    `json:"nomeItem"`
	Descricao          string    `json:"Descricao"`
	UnidadeMedida      string    `json:"unidadeMedida"`
	BoicoinsPorUnidade float64   `json:"boicoinsPorUnidade"`
}

func (id *ItemDoacao) TableName() string {
	return "boi_marronzinho.itens_doacao"
}
