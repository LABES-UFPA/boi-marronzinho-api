package repository

import (
    "boi-marronzinho-api/domain"
    "gorm.io/gorm"
)

type UserRepository interface {
    Repository[domain.Usuario]
    GetByEmail(email string) (*domain.Usuario, error)
}

type userRepository struct {
    Repository[domain.Usuario]
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        Repository: NewRepository[domain.Usuario](db),
        db: db,
    }
}

func (r *userRepository) GetByEmail(email string) (*domain.Usuario, error) {
    var user domain.Usuario
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
