package dto

import "github.com/google/uuid"

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	SaldoBoicoins float32   `json:"saldoBoicoins"`
	Token         *string   `json:"token"`
}

type UsuarioResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	SaldoBoicoins float32   `json:"saldoBoicoins"`
}
