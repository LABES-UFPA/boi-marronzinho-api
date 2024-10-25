package repository

import (
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrocaRepository interface {
    Repository[domain.Troca]
    BuscarTrocasCriadasAntesDe(data time.Time) ([]*domain.Troca, error)
    Deletar(trocaID uuid.UUID) error
    DeletarTrocasCriadasAntesDe(data time.Time) error
    ValidaTroca(idTroca uuid.UUID, validar bool) (*domain.Troca, error)

type trocaRepository struct {
    Repository[domain.Troca]
    db *gorm.DB
}

func NewTrocaRepository(db *gorm.DB) TrocaRepository {
    return &trocaRepository{
        Repository: NewRepository[domain.Troca](db),
        db:         db,
    }
}

func (r *trocaRepository) ValidaTroca(idTroca uuid.UUID, validar bool) (*domain.Troca, error) {
	var troca domain.Troca

	if err := r.db.Where("id = ?", idTroca).First(&troca).Error; err != nil {
		return nil, err
	}

	if troca.Status == "rejeitado" {
		return nil, errors.New("a troca já foi rejeitada")
	}
	if troca.Status == "validada" && validar {
		return nil, errors.New("a troca já foi validada")
	}

	novoStatus := "rejeitado"
	if validar {
		novoStatus = "validada"
	}

	if err := r.db.Model(&troca).UpdateColumn("status", novoStatus).Error; err != nil {
		return nil, err
	}


	troca.Status = novoStatus

	return &troca, nil
}

func (repo *trocaRepository) BuscarTrocasCriadasAntesDe(data time.Time) ([]*domain.Troca, error) {
    var trocas []*domain.Troca
    err := repo.db.Where("data_doacao < ?", data).Find(&trocas).Error
    return trocas, err
}

func (repo *trocaRepository) Deletar(trocaID uuid.UUID) error {
    return repo.db.Delete(&domain.Troca{}, trocaID).Error
}

func (repo *trocaRepository) DeletarTrocasCriadasAntesDe(data time.Time) error {
    return repo.db.Where("data_doacao < ? AND status = ?", data, "pendente").Delete(&domain.Troca{}).Error
}
