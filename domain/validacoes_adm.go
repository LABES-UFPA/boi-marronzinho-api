package domain

import (
	"time"

	"github.com/google/uuid"
)

type ValidacoesAdm struct {
	ID            uuid.UUID `json:"id"`
	Administador  uuid.UUID `json:"administador`
	Acao          string    `json:"acao"`
	DataValidacao time.Time `json:"datavalidacao"`
}
