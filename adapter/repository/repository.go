package repository

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Repository[T any] interface {
    Create(entity *T) (*T, error)
    Update(entity *T) (*T, error)
    Delete(id uuid.UUID) error
    GetByID(id uuid.UUID) (*T, error)
    GetAll() ([]*T, error)
    WithTransaction(fn func(txRepo Repository[T]) error) error
}

type repository[T any] struct {
    db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
    return &repository[T]{db: db}
}

func (r *repository[T]) Create(entity *T) (*T, error) {
    if err := r.db.Create(entity).Error; err != nil {
        return nil, err
    }
    return entity, nil
}

func (r *repository[T]) Update(entity *T) (*T, error) {
    if err := r.db.Save(entity).Error; err != nil {
        return nil, err
    }
    return entity, nil
}

func (r *repository[T]) Delete(id uuid.UUID) error {
    var entity T
    if err := r.db.Delete(&entity, id).Error; err != nil {
        return err
    }
    return nil
}

func (r *repository[T]) GetByID(id uuid.UUID) (*T, error) {
    var entity T
    if err := r.db.First(&entity, id).Error; err != nil {
        return nil, err
    }
    return &entity, nil
}

func (r *repository[T]) GetAll() ([]*T, error) {
    var entities []*T
    if err := r.db.Find(&entities).Error; err != nil {
        return nil, err
    }
    return entities, nil
}

func (r *repository[T]) WithTransaction(fn func(txRepo Repository[T]) error) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        txRepo := &repository[T]{db: tx}
        return fn(txRepo)
    })
}
