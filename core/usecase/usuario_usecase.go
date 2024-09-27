package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsuarioUseCase struct {
	userRepo repository.UserRepository
}

func NewUsuarioUseCase(userRepo repository.UserRepository) *UsuarioUseCase {
	return &UsuarioUseCase{userRepo: userRepo}
}

func (uc *UsuarioUseCase) CreateUser(usuarioRequest *domain.Usuario) (*domain.Usuario, error) {
	existingUser, err := uc.userRepo.GetByEmail(usuarioRequest.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	// hashedPassword, err := uc.HashPassword(usuarioRequest.Password)
	// if err != nil {
	// 	return nil, errors.New("erro ao gerar hash da senha")
	// }

	user := &domain.Usuario{
		ID:              uuid.New(),
		Nome:            usuarioRequest.Nome,
		Email:           usuarioRequest.Email,
		TipoUsuario:     usuarioRequest.TipoUsuario,
		IdiomaPreferido: "pt",
		CriadoEm:        time.Now(),
	}

	// err = user.Validate()
	// if err != nil {
	// 	return nil, err
	// }

	if err = uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
