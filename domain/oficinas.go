package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Oficinas struct {
	ID                  uuid.UUID `json:"id" gorm:"primaryKey"`
	Nome                string    `json:"nome" validate:"required"`
	Descricao           string    `json:"descricao" validate:"required"`
	PrecoBoicoins       float64   `json:"precoBoicoins" validate:"gt=0"`
	PrecoReal           float64   `json:"precoReal" validate:"gt=0"`
	DataEvento          time.Time `json:"dataEvento" validate:"required"`
	LimiteParticipantes int       `json:"limiteParticipantes" validate:"gt=0"`
	ParticipantesAtual  int       `json:"participantesAtual"`
	PontoMapaID         uuid.UUID `json:"pontoMapaId" validate:"required"`
	CriadoEm            time.Time `json:"-"`
}

func (o *Oficinas) TableName() string {
	return "boi_marronzinho.oficinas"
}

func (o *Oficinas) Validate() error {
	validate := validator.New()
	return validate.Struct(o)
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
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UsuarioID uuid.UUID `json:"usuarioId" gorm:"not null"`
	OficinaID uuid.UUID `json:"oficinaId" gorm:"not null"`
	Codigo    string    `json:"codigo" gorm:"size:100"`
	QRCode    []byte    `json:"qrcode" gorm:"type:bytea"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	Validado  bool      `json:"validado"`
}

func (tf *TicketOficina) TableName() string {
	return "boi_marronzinho.ticket_oficina"
}
