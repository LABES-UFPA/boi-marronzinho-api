package repository

import (
	"boi-marronzinho-api/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TrocaRepository interface {
	Repository[domain.Troca]
	BuscarTrocasCriadasAntesDe(data time.Time) ([]*domain.Troca, error)
	Deletar(trocaID uuid.UUID) error
	DeletarTrocasCriadasAntesDe(data time.Time) error
	ValidaTroca(idTroca uuid.UUID, validar bool) (*domain.Troca, error)
	
}

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

	logrus.Infof("Iniciando validação da troca com ID: %s", idTroca)

	if err := r.db.Where("id = ?", idTroca).First(&troca).Error; err != nil {
		logrus.Errorf("Erro ao buscar troca com ID %s: %v", idTroca, err)
		return nil, err
	}

	if troca.Status == "rejeitado" {
		logrus.Warnf("Tentativa de validar troca já rejeitada com ID: %s", idTroca)
		return nil, errors.New("a troca já foi rejeitada")
	}
	if troca.Status == "validada" && validar {
		logrus.Warnf("Tentativa de validar troca já validada com ID: %s", idTroca)
		return nil, errors.New("a troca já foi validada")
	}

	novoStatus := "rejeitado"
	if validar {
		novoStatus = "validada"
	}

	if err := r.db.Model(&troca).UpdateColumn("status", novoStatus).Error; err != nil {
		logrus.Errorf("Erro ao atualizar status da troca com ID %s: %v", idTroca, err)
		return nil, err
	}

	troca.Status = novoStatus
	logrus.Infof("Status da troca com ID %s atualizado para: %s", idTroca, novoStatus)

	return &troca, nil
}

func (repo *trocaRepository) BuscarTrocasCriadasAntesDe(data time.Time) ([]*domain.Troca, error) {
	logrus.Infof("Buscando trocas criadas antes de %v", data)
	var trocas []*domain.Troca
	err := repo.db.Where("data_doacao < ?", data).Find(&trocas).Error
	if err != nil {
		logrus.Errorf("Erro ao buscar trocas criadas antes de %v: %v", data, err)
	}
	return trocas, err
}

func (repo *trocaRepository) Deletar(trocaID uuid.UUID) error {
	logrus.Infof("Deletando troca com ID: %s", trocaID)
	if err := repo.db.Delete(&domain.Troca{}, trocaID).Error; err != nil {
		logrus.Errorf("Erro ao deletar troca com ID %s: %v", trocaID, err)
		return err
	}
	logrus.Infof("Troca com ID %s deletada com sucesso", trocaID)
	return nil
}

// func (repo *trocaRepository) validacaoAdm(admID *domain.ValidacoesAdm) error {
// 	logrus.Infof("Validando a troca pelo administrador com ID: %s", admID.ID)
// 	if err := repo.db.Create(admID).Error; err != nil {
// 		logrus.Errorf("Erro ao validar troca pelo administrador com ID %s: %v", admID.ID, err)
// 		return err
// 	}
// 	logrus.Infof("Validação da troca pelo administrador com ID %s realizada com sucesso", admID.ID)
// 	return nil
// }

func (repo *trocaRepository) DeletarTrocasCriadasAntesDe(data time.Time) error {
	logrus.Infof("Deletando trocas criadas antes de %v com status 'pendente'", data)

	result := repo.db.Where("data_troca < ? AND status = ?", data, "pendente").Delete(&domain.Troca{})
	if err := result.Error; err != nil {
		logrus.Errorf("Erro ao deletar trocas criadas antes de %v: %v", data, err)
		return err
	}

	quantidadeDeletada := result.RowsAffected
	logrus.Infof("Deletadas %d trocas criadas antes de %v com status 'pendente'", quantidadeDeletada, data)

	return nil
}
