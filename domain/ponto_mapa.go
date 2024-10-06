package domain

import "github.com/google/uuid"

type PontosMapa struct {
	ID        uuid.UUID `json:"id" gorm:"PrimaryKey"`
	Nome      string    `json:"nome"`
	Descricao string    `json:"descricao"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	ImagemURL string    `json:"imagemURL"`
	CriadoEm  string    `json:"-"`
}

func (pm *PontosMapa) TableName() string {
	return "boi_marronzinho.pontos_mapa"
}
