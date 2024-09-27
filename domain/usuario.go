package domain

import (
	"github.com/google/uuid"
	"time"
)

type Usuario struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	Nome            string    `json:"nome"`
	Email           string    `json:"email"`
	TipoUsuario     string    `json:"tipo_usuario"`
	IdiomaPreferido string    `json:"idioma_preferido"`
	CriadoEm        time.Time `json:"criado_em"`
}

func (d *Usuario) TableName() string {
	return "boi_marronzinho.usuarios"
}
