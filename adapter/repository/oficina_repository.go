package repository

import (
	"boi-marronzinho-api/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OficinaRepository interface {
	Repository[domain.Oficinas]
	InscreverParticipante(oficinaID uuid.UUID, usuario *domain.Usuario, pagamentoEmBoicoins bool) error
}

type oficinaRepository struct {
	Repository[domain.Oficinas]
	db *gorm.DB
}

func NewOficinaRepository(db *gorm.DB) OficinaRepository {
	return &oficinaRepository{
		Repository: NewRepository[domain.Oficinas](db),
		db:         db,
	}
}

func (r *oficinaRepository) InscreverParticipante(oficinaID uuid.UUID, usuario *domain.Usuario, pagamentoEmBoicoins bool) error {
	var oficina domain.Oficinas

	if err := r.db.First(&oficina, "id = ?", oficinaID).Error; err != nil {
		return err
	}

	if pagamentoEmBoicoins {
		transacao := domain.BoicoinsTransacoes{
			ID:            uuid.New(),
			UsuarioID:     usuario.ID,
			Quantidade:    -float64(oficina.PrecoBoicoins),
			TipoTransacao: "inscricao_oficina",
			Descricao:     "Inscrição na oficina " + oficina.Nome,
			DataTransacao: time.Now(),
			OficinaID:     &oficina.ID,
			DoacaoID:      nil,
			PedidoID:      nil,
		}

		if err := r.db.Create(&transacao).Error; err != nil {
			return err
		}
	}

	oficina.ParticipantesAtual++
	if err := r.db.Save(&oficina).Error; err != nil {
		return err
	}

	inscricao := domain.TicketOficina{
		ID:            uuid.New(),
		UsuarioID:     usuario.ID,
		OficinaID:     oficina.ID,
		Codigo:        uuid.New().String(),
		DataInscricao: time.Now(),
	}

	if err := r.db.Create(&inscricao).Error; err != nil {
		return err
	}

	po := &domain.ParticipanteOficina{
		ID:            uuid.New(),
		UsuarioID:     usuario.ID,
		OficinaID:     oficina.ID,
		DataInscricao: time.Now(),
	}

	if err := r.db.Create(&po).Error; err != nil {
		return err
	}

	return nil
}
