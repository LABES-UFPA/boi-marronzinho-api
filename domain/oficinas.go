package domain

import (
	"time"

	"github.com/google/uuid"
)

type Oficinas struct {
	ID                  uuid.UUID `json:"id" gorm:"primaryKey"`
	Nome                string    `json:"nome"`
	Descricao           string    `json:"descricao"`
	PrecoBoicoins       float64   `json:"precoBoicoins"`
	PrecoReal           float64   `json:"precoReal"`
	DataEvento          time.Time `json:"dataEvento"`
	LimiteParticipantes int       `json:"limiteParticipantes"`
	ParticipantesAtual  int       `json:"participantesAtual"`
	PontoMapaID         uuid.UUID `json:"pontoMapaId"`
	CriadoEm            time.Time `json:"-"`
}

func (o *Oficinas) TableName() string {
	return "boi_marronzinho.oficinas"
}

type ParticipanteOficina struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID `json:"usuarioId" gorm:"not null"`
	OficinaID     uuid.UUID `json:"oficinaId" gorm:"not null"`
	DataInscricao time.Time `json:"dataInscricao" gorm:"default:CURRENT_TIMESTAMP"`
}

func (po *ParticipanteOficina) TableName() string {
	return "boi_marronzinho.participantes_oficinas"
}

type TicketOficina struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID     uuid.UUID `json:"usuarioId" gorm:"not null"`
	OficinaID     uuid.UUID `json:"oficinaId" gorm:"not null"`
	Codigo        string    `json:"codigo" gorm:"size:100"` // Código de validação ou QR Code
	DataInscricao time.Time `json:"dataInscricao" gorm:"default:CURRENT_TIMESTAMP"`
}

func (tf *TicketOficina) TableName() string {
	return "boi_marronzinho.ticket_oficina"
}
