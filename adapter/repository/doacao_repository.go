package repository

import (
	"boi-marronzinho-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoacaoRepository interface {
	AdicionaDoacao(doacao *domain.Doacoes) (*domain.Doacoes, error)
	AtualizaItemDoacao(itemDoacao *domain.ItemDoacao) (*domain.ItemDoacao, error)
	DeletaItemDoacao(id uuid.UUID) error
	CriaItemDoacao(itemDoacao *domain.ItemDoacao) (*domain.ItemDoacao, error)
	CapturaTodosItensDoacao() ([]*domain.ItemDoacao, error)
	CapturaItemDoacao(id uuid.UUID) (*domain.ItemDoacao, error)
}

type doacaoRepository struct {
	db *gorm.DB
}

func NewDoacaoRepository(db *gorm.DB) DoacaoRepository {
	return &doacaoRepository{db: db}
}

func (r *doacaoRepository) AdicionaDoacao(doacao *domain.Doacoes) (*domain.Doacoes, error) {
	if err := r.db.Create(doacao).Error; err != nil {
		return nil, err
	}
	return doacao, nil
}

func (r *doacaoRepository) AtualizaItemDoacao(itemDoacao *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	if err := r.db.Save(itemDoacao).Error; err != nil {
		return nil, err
	}
	return itemDoacao, nil
}

func (r *doacaoRepository) DeletaItemDoacao(id uuid.UUID) error {
	if err := r.db.Delete(&domain.ItemDoacao{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *doacaoRepository) CriaItemDoacao(itemDoacao *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	if err := r.db.Create(itemDoacao).Error; err != nil {
		return nil, err
	}
	return itemDoacao, nil
}

func (r *doacaoRepository) CapturaTodosItensDoacao() ([]*domain.ItemDoacao, error) {
	var itens []*domain.ItemDoacao
	if err := r.db.Find(&itens).Error; err != nil {
		return nil, err
	}
	return itens, nil
}

func (r *doacaoRepository) CapturaItemDoacao(id uuid.UUID) (*domain.ItemDoacao, error) {
	var itemDoacao domain.ItemDoacao
	if err := r.db.First(&itemDoacao, id).Error; err != nil {
		return nil, err
	}
	return &itemDoacao, nil
}
