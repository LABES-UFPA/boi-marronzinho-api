package dto

import "boi-marronzinho-api/domain"

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Usuario *domain.Usuario
	Token   *string `json:"token"`
}
