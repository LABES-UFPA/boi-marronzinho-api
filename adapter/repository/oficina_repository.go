package repository

import (
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	"boi-marronzinho-api/minio"
	"bytes"
	"errors"
	"fmt"   
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type OficinaRepository interface {
	Repository[domain.Oficinas]
	InscreverParticipante(oficinaID uuid.UUID, usuario *domain.Usuario, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error)
	BuscarUsuarioPorCodigo(codigoTicket string) (*domain.Usuario, error)
	GetTicketsByUsuarioID(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error)
	ValidaVoucher(codigo *string) (*domain.TicketOficina, error)
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

func (r *oficinaRepository) InscreverParticipante(oficinaID uuid.UUID, usuario *domain.Usuario, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error) {
    var oficina domain.Oficinas

    if err := r.db.First(&oficina, "id = ?", oficinaID).Error; err != nil {
        return nil, err
    }

    var participanteExistente domain.ParticipanteOficina
    if err := r.db.Where("usuario_id = ? AND oficina_id = ?", usuario.ID, oficinaID).First(&participanteExistente).Error; err == nil {
        return nil, errors.New("usuário já está inscrito nesta oficina")
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
        }

        if err := r.db.Create(&transacao).Error; err != nil {
            return nil, err
        }
    }

    oficina.ParticipantesAtual++
    if err := r.db.Save(&oficina).Error; err != nil {
        return nil, err
    }

    codigoTicket := uuid.New().String()
    png, err := qrcode.Encode(codigoTicket, qrcode.Medium, 256)
    if err != nil {
        return nil, err
    }
    qrCodeFileName := fmt.Sprintf("%s.png", codigoTicket)

    qrCodeURL, err := minio.UploadFile(bytes.NewReader(png), qrCodeFileName, "qrcodes")
    if err != nil {
        return nil, err
    }

    inscricao := domain.TicketOficina{
        ID:        uuid.New(),
        UsuarioID: usuario.ID,
        OficinaID: oficina.ID,
        Codigo:    codigoTicket,
        QRCode:    qrCodeURL,
        CreatedAt: time.Now(),
    }

    if err := r.db.Create(&inscricao).Error; err != nil {
        return nil, err
    }

    po := &domain.ParticipanteOficina{
        ID:            uuid.New(),
        UsuarioID:     usuario.ID,
        OficinaID:     oficina.ID,
        DataInscricao: time.Now(),
    }

    if err := r.db.Create(&po).Error; err != nil {
        return nil, err
    }

    return po, nil
}


func (r *oficinaRepository) BuscarUsuarioPorCodigo(codigoTicket string) (*domain.Usuario, error) {
	var ticket domain.TicketOficina

	if err := r.db.Where("codigo = ?", codigoTicket).First(&ticket).Error; err != nil {
		return nil, err
	}

	var usuario domain.Usuario
	if err := r.db.Where("id = ?", ticket.UsuarioID).First(&usuario).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (r *oficinaRepository) GetTicketsByUsuarioID(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error) {
	var results []dto.VoucherResponseDTO

	err := r.db.Table("ticket_oficina").
		Select("ticket_oficina.id, ticket_oficina.usuario_id, ticket_oficina.oficina_id, ticket_oficina.qr_code, ticket_oficina.created_at, ticket_oficina.validado, oficinas.nome AS nome_oficina, oficinas.descricao AS descricao").
		Joins("JOIN oficinas ON ticket_oficina.oficina_id = oficinas.id").
		Where("ticket_oficina.usuario_id = ?", usuarioID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *oficinaRepository) ValidaVoucher(codigoVoucher *string) (*domain.TicketOficina, error) {
	var ticket domain.TicketOficina

	if err := r.db.Where("codigo = ?", *codigoVoucher).First(&ticket).Error; err != nil {
		return nil, err
	}

	if ticket.Validado {
		return nil, errors.New("o ticket já foi validado")
	}

	ticket.Validado = true
	if err := r.db.Save(&ticket).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}
