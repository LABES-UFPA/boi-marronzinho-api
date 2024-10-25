package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Evento struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	Nome         string    `json:"nome" validate:"required"`
	Descricao    string    `json:"descricao" validate:"required"`
	DataEvento   time.Time `json:"dataEvento" validate:"required"`
	LinkEndereco string    `json:"linkEndereco" validate:"required"`
	ImagemUrl    string    `json:"imagemUrl"`
	CriadoEm     time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
	AtualizadoEm time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
}

func (e *Evento) TableName() string {
	return "boi_marronzinho.evento"
}


func (e *Evento) Validate() error {
	validate := validator.New()
	return validate.Struct(e)
}
