package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.Usuario) error
	GetByID(id int) (*domain.Usuario, error)
	GetByEmail(email string) (*domain.Usuario, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(usuario *domain.Usuario) error {
	logrus.WithFields(logrus.Fields{
		"user":  usuario.FirstName + " " + usuario.LastName,
		"email": usuario.Email,
	}).Info("Create new user")

	return r.db.Create(usuario).Error
}

func (r *userRepository) GetByID(id int) (*domain.Usuario, error) {
	var user domain.Usuario
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.Usuario, error) {
	var user domain.Usuario
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
