package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/auth"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuarioRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("falha ao criptografar a senha")
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

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	if _, err = uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Login(email, password string) (*dto.LoginResponseDTO, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciais inválidas")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	token, err := auth.GenerateJWT(user)
	if err != nil {
		return nil, errors.New("não foi possível gerar o token")
	}

	userResponse := &dto.LoginResponseDTO{
		Token:         &token,
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		SaldoBoicoins: user.SaldoBoicoins,
	}

	return userResponse, nil
}

func (uc *UserUseCase) GetUserByID(id uuid.UUID) (*dto.UsuarioResponseDTO, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	return &dto.UsuarioResponseDTO{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		SaldoBoicoins: user.SaldoBoicoins,
	}, nil
}

func (uc *UserUseCase) UpdateUser(id uuid.UUID, updateData *domain.Usuario) (*domain.Usuario, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.TipoUsuario != "" {
		user.TipoUsuario = updateData.TipoUsuario
	}
	if updateData.IdiomaPreferido != "" {
		user.IdiomaPreferido = updateData.IdiomaPreferido
	}
	user.UpdatedAt = time.Now()

	if _, err = uc.userRepo.Update(user); err != nil {
		return nil, errors.New("falha ao atualizar o usuário")
	}

	return user, nil
}

func (uc *UserUseCase) DeleteUser(id uuid.UUID) error {
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuário não encontrado")
		}
		return err
	}

	if err := uc.userRepo.Delete(id); err != nil {
		return errors.New("falha ao deletar o usuário")
	}

	return nil
}


func (uc *UserUseCase) GetExtrato(id uuid.UUID) ([]*domain.BoicoinsTransacoes, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	extrato, err := uc.userRepo.GetExtrato(user.ID)
	if err!= nil {
        return nil, err
    }

	return extrato, nil
}