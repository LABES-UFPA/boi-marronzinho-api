package domain

import (
	"time"
	"github.com/google/uuid"
)

type Usuario struct {
	ID             uuid.UUID      `json:"id"`                        
	Nome           string         `json:"nome"`                    
	Email          string         `json:"email"`                   
	TipoUsuario    string         `json:"tipo_usuario"`     
	IdiomaPreferido string        `json:"idioma_preferido"`
	CriadoEm       time.Time      `json:"criado_em"`
}