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
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type OficinaRepository interface {
	Repository[domain.Oficinas]
	InscreverParticipante(oficinaID uuid.UUID, usuario *domain.Usuario, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error)
	BuscarUsuarioPorCodigo(codigoTicket string) (*domain.Usuario, error)
	GetTicketsByUsuarioID(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error)
	ValidaVoucher(codigo *string) (*domain.TicketOficina, error)
	DeleteOficina(id uuid.UUID) error
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
	logrus.Infof("Iniciando inscrição de participante para oficina ID: %s, Usuário ID: %s", oficinaID, usuario.ID)

	var oficina domain.Oficinas
	if err := r.db.First(&oficina, "id = ?", oficinaID).Error; err != nil {
		logrus.Errorf("Erro ao buscar oficina com ID %s: %v", oficinaID, err)
		return nil, err
	}

	var participanteExistente domain.ParticipanteOficina
	if err := r.db.Where("usuario_id = ? AND oficina_id = ?", usuario.ID, oficinaID).First(&participanteExistente).Error; err == nil {
		logrus.Warnf("Usuário com ID %s já está inscrito na oficina ID %s", usuario.ID, oficinaID)
		return nil, errors.New("usuário já está inscrito nesta oficina")
	}

	if pagamentoEmBoicoins {
		logrus.Infof("Processando pagamento em Boicoins para usuário ID: %s", usuario.ID)
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
			logrus.Errorf("Erro ao criar transação de Boicoins para usuário ID %s: %v", usuario.ID, err)
			return nil, err
		}
		logrus.Infof("Pagamento em Boicoins processado com sucesso para usuário ID: %s", usuario.ID)
	}

	oficina.ParticipantesAtual++
	if err := r.db.Save(&oficina).Error; err != nil {
		logrus.Errorf("Erro ao atualizar contagem de participantes para oficina ID %s: %v", oficinaID, err)
		return nil, err
	}

	codigoTicket := uuid.New().String()
	png, err := qrcode.Encode(codigoTicket, qrcode.Medium, 256)
	if err != nil {
		logrus.Errorf("Erro ao gerar QR Code para ticket: %v", err)
		return nil, err
	}

	qrCodeFileName := fmt.Sprintf("%s.png", codigoTicket)
	qrCodeURL, err := minio.UploadFile(bytes.NewReader(png), qrCodeFileName, "qrcodes")
	if err != nil {
		logrus.Errorf("Erro ao fazer upload do QR Code para ticket %s: %v", codigoTicket, err)
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
		logrus.Errorf("Erro ao criar ticket de oficina para usuário ID %s: %v", usuario.ID, err)
		return nil, err
	}

	po := &domain.ParticipanteOficina{
		ID:            uuid.New(),
		UsuarioID:     usuario.ID,
		OficinaID:     oficina.ID,
		DataInscricao: time.Now(),
	}

	if err := r.db.Create(&po).Error; err != nil {
		logrus.Errorf("Erro ao registrar participante de oficina para usuário ID %s: %v", usuario.ID, err)
		return nil, err
	}

	logrus.Infof("Inscrição de participante concluída com sucesso para oficina ID: %s, Usuário ID: %s", oficinaID, usuario.ID)
	return po, nil
}

func (r *oficinaRepository) BuscarUsuarioPorCodigo(codigoTicket string) (*domain.Usuario, error) {
	logrus.Infof("Buscando usuário pelo código do ticket: %s", codigoTicket)
	var ticket domain.TicketOficina

	if err := r.db.Where("codigo = ?", codigoTicket).First(&ticket).Error; err != nil {
		logrus.Errorf("Erro ao buscar ticket com código %s: %v", codigoTicket, err)
		return nil, err
	}

	var usuario domain.Usuario
	if err := r.db.Where("id = ?", ticket.UsuarioID).First(&usuario).Error; err != nil {
		logrus.Errorf("Erro ao buscar usuário com ID %s: %v", ticket.UsuarioID, err)
		return nil, err
	}

	logrus.Infof("Usuário encontrado para código do ticket %s: Usuário ID: %s", codigoTicket, usuario.ID)
	return &usuario, nil
}

func (r *oficinaRepository) GetTicketsByUsuarioID(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error) {
	logrus.Infof("Buscando tickets para usuário com ID: %s", usuarioID)
	var results []dto.VoucherResponseDTO

	err := r.db.Table("ticket_oficina").
		Select("ticket_oficina.id, ticket_oficina.usuario_id, ticket_oficina.oficina_id, ticket_oficina.qr_code, ticket_oficina.created_at, ticket_oficina.validado, oficinas.nome AS nome_oficina, oficinas.descricao AS descricao").
		Joins("JOIN oficinas ON ticket_oficina.oficina_id = oficinas.id").
		Where("ticket_oficina.usuario_id = ?", usuarioID).
		Scan(&results).Error

	if err != nil {
		logrus.Errorf("Erro ao buscar tickets para usuário com ID %s: %v", usuarioID, err)
		return nil, err
	}

	logrus.Infof("Encontrados %d tickets para usuário com ID: %s", len(results), usuarioID)
	return results, nil
}

func (r *oficinaRepository) ValidaVoucher(codigoVoucher *string) (*domain.TicketOficina, error) {
	logrus.Infof("Validando voucher com código: %s", *codigoVoucher)
	var ticket domain.TicketOficina

	if err := r.db.Where("codigo = ?", *codigoVoucher).First(&ticket).Error; err != nil {
		logrus.Errorf("Erro ao buscar voucher com código %s: %v", *codigoVoucher, err)
		return nil, err
	}

	if ticket.Validado {
		logrus.Warnf("Tentativa de validar voucher já validado com código: %s", *codigoVoucher)
		return nil, errors.New("o ticket já foi validado")
	}

	ticket.Validado = true
	if err := r.db.Save(&ticket).Error; err != nil {
		logrus.Errorf("Erro ao salvar voucher validado com código %s: %v", *codigoVoucher, err)
		return nil, err
	}

	logrus.Infof("Voucher validado com sucesso: %s", *codigoVoucher)
	return &ticket, nil
}


func (r *oficinaRepository) DeleteOficina(id uuid.UUID) error {
	logrus.Infof("Iniciando exclusão da oficina ID: %s", id)

	var oficina domain.Oficinas
	if err := r.db.First(&oficina, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("Oficina com ID %s não encontrada: %v", id, err)
			return errors.New("oficina não encontrada")
		}
		return err
	}

	var inscritoExistente bool
	if err := r.db.Model(&domain.ParticipanteOficina{}).Where("oficina_id = ?", id).Select("count(*) > 0").Find(&inscritoExistente).Error; err != nil {
		logrus.Errorf("Erro ao verificar inscrições para oficina ID %s: %v", id, err)
		return err
	}

	if inscritoExistente {
		logrus.Warnf("Não é possível excluir a oficina ID %s pois existem participantes inscritos.", id)
		return errors.New("não é possível excluir uma oficina com participantes inscritos")
	}

	if err := r.db.Delete(&oficina).Error; err != nil {
		logrus.Errorf("Erro ao deletar oficina ID %s: %v", id, err)
		return err
	}

	logrus.Infof("Oficina ID %s deletada com sucesso", id)
	return nil
}
