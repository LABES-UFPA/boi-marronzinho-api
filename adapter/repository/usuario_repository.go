package repository

import (
	"boi-marronzinho-api/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[domain.Usuario]
	GetByEmail(email string) (*domain.Usuario, error)
	AtualizarSaldo(usuarioID uuid.UUID, boicoinsRecebidos float64) error
	GetExtrato(usuarioID uuid.UUID) ([]*domain.BoicoinsTransacoes, error)
	GetByFullName(name string) ([]*domain.Usuario, error)
}

type userRepository struct {
	Repository[domain.Usuario]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		Repository: NewRepository[domain.Usuario](db),
		db:         db,
	}
}

func (ur *userRepository) GetByEmail(email string) (*domain.Usuario, error) {
	var user domain.Usuario
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) AtualizarSaldo(usuarioID uuid.UUID, boicoinsRecebidos float64) error {
	result := ur.db.Model(&domain.Usuario{}).Where("id = ?", usuarioID).
		Update("saldo_boicoins", gorm.Expr("saldo_boicoins + ?", boicoinsRecebidos))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuário não encontrado ou saldo não atualizado")
	}
	return nil
}

func (ur *userRepository) GetExtrato(usuarioID uuid.UUID) ([]*domain.BoicoinsTransacoes, error) {
	var transacoes []*domain.BoicoinsTransacoes
	if err := ur.db.Where("usuario_id = ?", usuarioID).Order("data_transacao DESC").Find(&transacoes).Error; err != nil {
		return nil, err
	}
	return transacoes, nil
}

func (ur *userRepository) GetByFullName(name string) ([]*domain.Usuario, error) {
	var users []*domain.Usuario
	query := "%" + name + "%"
	if err := ur.db.Where("CONCAT(first_name, ' ', last_name) LIKE ?", query).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
