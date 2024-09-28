package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usuario struct {
	ID                   uuid.UUID      `json:"id" gorm:"primaryKey"`
	FirstName            string         `json:"first_name" validate:"required"`
	LastName             string         `json:"last_name" validate:"required"`
	Email                string         `json:"email" validate:"required"`
	TipoUsuario          string         `json:"tipoUsuario" validate:"required"`
	IdiomaPreferido      string         `json:"idiomaPreferido" validate:"required"`
	Password             string         `json:"password" validate:"required" gorm:"-"`
	PasswordHash         string         `json:"-" validate:"required,min=1,max=255"`
	LastLogin            *time.Time     `json:"_"`
	PasswordResetToken   *string        `json:"-"`
	PasswordResetExpires *time.Time     `json:"-"`
	CreatedAt            time.Time      `json:"-"`
	UpdatedAt            time.Time      `json:"-"`
	DeletedAt            gorm.DeletedAt `json:"-"`
}

func (u *Usuario) TableName() string {
	return "boi_marronzinho.usuarios"
}

func (u *Usuario) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
