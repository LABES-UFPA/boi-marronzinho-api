package repository

import (
	"boi-marronzinho-api/domain"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Buscando usuário com email: %s", email)
	var user domain.Usuario
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		logrus.Errorf("Erro ao buscar usuário com email %s: %v", email, err)
		return nil, err
	}
	logrus.Infof("Usuário encontrado com email: %s, ID: %s", email, user.ID)
	return &user, nil
}

func (ur *userRepository) AtualizarSaldo(usuarioID uuid.UUID, boicoinsRecebidos float64) error {
	logrus.Infof("Atualizando saldo do usuário com ID: %s, adicionando %f Boicoins", usuarioID, boicoinsRecebidos)
	result := ur.db.Model(&domain.Usuario{}).Where("id = ?", usuarioID).
		Update("saldo_boicoins", gorm.Expr("saldo_boicoins + ?", boicoinsRecebidos))
	if result.Error != nil {
		logrus.Errorf("Erro ao atualizar saldo do usuário com ID %s: %v", usuarioID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logrus.Warnf("Usuário com ID %s não encontrado ou saldo não atualizado", usuarioID)
		return errors.New("usuário não encontrado ou saldo não atualizado")
	}
	logrus.Infof("Saldo atualizado com sucesso para o usuário com ID: %s", usuarioID)
	return nil
}

func (ur *userRepository) GetExtrato(usuarioID uuid.UUID) ([]*domain.BoicoinsTransacoes, error) {
	logrus.Infof("Buscando extrato de Boicoins para o usuário com ID: %s", usuarioID)
	var transacoes []*domain.BoicoinsTransacoes
	if err := ur.db.Where("usuario_id = ?", usuarioID).Order("data_transacao DESC").Find(&transacoes).Error; err != nil {
		logrus.Errorf("Erro ao buscar extrato para o usuário com ID %s: %v", usuarioID, err)
		return nil, err
	}
	logrus.Infof("Extrato de Boicoins encontrado para o usuário com ID: %s", usuarioID)
	return transacoes, nil
}

func (ur *userRepository) GetByFullName(name string) ([]*domain.Usuario, error) {
	logrus.Infof("Buscando usuários com nome completo que contém: %s", name)
	var users []*domain.Usuario
	query := "%" + name + "%"
	if err := ur.db.Where("CONCAT(first_name, ' ', last_name) LIKE ?", query).Find(&users).Error; err != nil {
		logrus.Errorf("Erro ao buscar usuários com nome que contém %s: %v", name, err)
		return nil, err
	}
	logrus.Infof("Encontrados %d usuários com nome completo que contém: %s", len(users), name)
	return users, nil
}
