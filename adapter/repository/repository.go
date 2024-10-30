package repository

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*T, error)
	GetAll() ([]*T, error)
	BatchCreate(entities []*T) error
	BatchUpdate(entities []*T) error
	BatchDelete(ids []uuid.UUID) error
}

type repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &repository[T]{db: db}
}

func (r *repository[T]) Create(entity *T) (*T, error) {
	logrus.Infof("Criando nova entidade: %+v", entity)
	if err := r.db.Create(entity).Error; err != nil {
		logrus.Errorf("Erro ao criar entidade: %v", err)
		return nil, err
	}
	logrus.Infof("Entidade criada com sucesso: %+v", entity)
	return entity, nil
}

func (r *repository[T]) Update(entity *T) (*T, error) {
	logrus.Infof("Atualizando entidade: %+v", entity)
	if err := r.db.Save(entity).Error; err != nil {
		logrus.Errorf("Erro ao atualizar entidade: %v", err)
		return nil, err
	}
	logrus.Infof("Entidade atualizada com sucesso: %+v", entity)
	return entity, nil
}

func (r *repository[T]) Delete(id uuid.UUID) error {
	logrus.Infof("Deletando entidade com ID: %s", id)
	var entity T
	if err := r.db.Delete(&entity, id).Error; err != nil {
		logrus.Errorf("Erro ao deletar entidade com ID %s: %v", id, err)
		return err
	}
	logrus.Infof("Entidade deletada com sucesso com ID: %s", id)
	return nil
}

func (r *repository[T]) GetByID(id uuid.UUID) (*T, error) {
	logrus.Infof("Buscando entidade com ID: %s", id)
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		logrus.Errorf("Erro ao buscar entidade com ID %s: %v", id, err)
		return nil, err
	}
	logrus.Infof("Entidade encontrada com ID: %s", id)
	return &entity, nil
}

func (r *repository[T]) GetAll() ([]*T, error) {
	logrus.Infof("Buscando todas as entidades")
	var entities []*T
	if err := r.db.Find(&entities).Error; err != nil {
		logrus.Errorf("Erro ao buscar todas as entidades: %v", err)
		return nil, err
	}
	logrus.Infof("Encontradas %d entidades", len(entities))
	return entities, nil
}

func (r *repository[T]) BatchCreate(entities []*T) error {
	logrus.Infof("Criando múltiplas entidades em batch")
	tx := r.db.Begin()
	for _, entity := range entities {
		if err := tx.Create(entity).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("Erro ao criar entidade em batch: %v", err)
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("Erro ao commitar transação de criação em batch: %v", err)
		return err
	}
	logrus.Infof("Criação em batch concluída com sucesso")
	return nil
}

func (r *repository[T]) BatchUpdate(entities []*T) error {
	logrus.Infof("Atualizando múltiplas entidades")
	for _, entity := range entities {
		if err := r.db.Save(entity).Error; err != nil {
			logrus.Errorf("Erro ao atualizar entidade: %v", err)
			return err
		}
	}
	logrus.Infof("Atualização em batch concluída com sucesso")
	return nil
}

func (r *repository[T]) BatchDelete(ids []uuid.UUID) error {
	logrus.Infof("Deletando múltiplas entidades")
	if err := r.db.Delete(new(T), ids).Error; err != nil {
		logrus.Errorf("Erro ao deletar múltiplas entidades: %v", err)
		return err
	}
	logrus.Infof("Deleção em batch concluída com sucesso")
	return nil
}
