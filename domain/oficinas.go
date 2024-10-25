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
	LimiteParticipantes *int      `json:"limiteParticipantes" validate:"gt=0"`
	ParticipantesAtual  int       `json:"participantesAtual"`
	ImagemUrl           string    `json:"imagemUrl"`
	LinkEndereco        string    `json:"linkEndereco"`
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
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UsuarioID uuid.UUID `gorm:"type:uuid"`
	OficinaID uuid.UUID `gorm:"type:uuid"`
	Codigo    string    `gorm:"unique"`
	QRCode    string    `gorm:"type:varchar(255)"`
	Validado  bool      `gorm:"boolean"`
	CreatedAt time.Time
}

func (tf *TicketOficina) TableName() string {
	return "boi_marronzinho.ticket_oficina"
}
