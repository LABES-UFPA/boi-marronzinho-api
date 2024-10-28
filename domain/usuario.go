package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usuario struct {
	ID                   uuid.UUID      `json:"id" gorm:"primaryKey"`
	FirstName            string         `json:"firstName" validate:"required"`
	LastName             string         `json:"lastName" validate:"required"`
	Email                string         `json:"email" validate:"required,email"`                                    // Validação para formato de email
	TipoUsuario          string         `json:"tipoUsuario" validate:"required,oneof=Cliente Administrador Gringo"` // Validação para tipos de usuário permitidos
	IdiomaPreferido      string         `json:"idiomaPreferido" validate:"required"`
	Password             string         `json:"password" validate:"required,min=6" gorm:"-"`
	PasswordHash         string         `json:"-" validate:"required,min=1,max=255"`
	SaldoBoicoins        float32        `json:"-" gorm:"default:0.00"`
	LastLogin            *time.Time     `json:"lastLogin,omitempty"`
	PasswordResetToken   *string        `json:"passwordResetToken,omitempty"`
	PasswordResetExpires *time.Time     `json:"passwordResetExpires,omitempty"`
	CreatedAt            time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

func (u *Usuario) TableName() string {
    return "boi_marronzinho.usuarios"
}

func (u *Usuario) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
