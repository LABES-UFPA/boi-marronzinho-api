package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/auth"
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUsuarioUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) CreateUser(usuarioRequest *domain.Usuario) (*domain.Usuario, error) {
	existingUser, err := uc.userRepo.GetByEmail(usuarioRequest.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuarioRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.Usuario{
		ID:              uuid.New(),
		FirstName:       usuarioRequest.FirstName,
		LastName:        usuarioRequest.LastName,
		Email:           usuarioRequest.Email,
		PasswordHash:    string(hashedPassword),
		TipoUsuario:     usuarioRequest.TipoUsuario,
		IdiomaPreferido: usuarioRequest.IdiomaPreferido,
		CreatedAt:       time.Now(),
	}

	if err = uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(user.Email, user.TipoUsuario)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}
