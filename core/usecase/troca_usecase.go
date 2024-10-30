package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type TrocaUseCase struct {
	trocaRepo     repository.TrocaRepository
	ItemTrocaRepo repository.Repository[domain.ItemTroca]
	userRepo      repository.UserRepository
	boicoinRepo   repository.BoicoinRepository
}

func NewTrocaUseCase(
	trocaRepo repository.TrocaRepository,
	ItemTrocaRepo repository.Repository[domain.ItemTroca],
	userRepo repository.UserRepository,
	boicoinRepo repository.BoicoinRepository,
) *TrocaUseCase {
	return &TrocaUseCase{
		trocaRepo:     trocaRepo,
		ItemTrocaRepo: ItemTrocaRepo,
		userRepo:      userRepo,
		boicoinRepo:   boicoinRepo,
	}
}

func (tuc *TrocaUseCase) RealizarTroca(trocaRequest *domain.Troca) (*domain.Troca, string, error) {
	_, err := tuc.ItemTrocaRepo.GetByID(trocaRequest.ItemTrocaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("item de doação não encontrado")
		}
		return nil, "", err
	}

	troca := &domain.Troca{
		ID:                uuid.New(),
		UsuarioID:         trocaRequest.UsuarioID,
		ItemTrocaID:       trocaRequest.ItemTrocaID,
		Quantidade:        trocaRequest.Quantidade,
		BoicoinsRecebidos: trocaRequest.BoicoinsRecebidos,
		DataTroca:         time.Now(),
		Status:            "pendente",
	}

	createdtroca, err := tuc.trocaRepo.Create(troca)
	if err != nil {
		return nil, "", err
	}

	qrCodeBase64, err := gerarQRCodeBase64(createdtroca)
	if err != nil {
		logrus.Error("Erro ao gerar QR Code: ", err)
		return nil, "", errors.New("erro ao gerar QR Code")
	}

	//go func() {
	//	if err := tuc.notificarAdministrador(createdtroca); err != nil {
	//		logrus.Info("Erro ao notificar administrador" + err.Error())
	//	}
	//}()

	return createdtroca, qrCodeBase64, nil
}

func gerarQRCodeBase64(troca *domain.Troca) (string, error) {
	qrData := troca.ID.String()
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(png)
	return qrCodeBase64, nil
}

func (tuc *TrocaUseCase) ValidaTroca(trocaID uuid.UUID, validar bool) (*domain.Troca, error) {
	_, err := tuc.trocaRepo.GetByID(trocaID)
	if err != nil {
		return nil, err
	}

	t, err := tuc.trocaRepo.ValidaTroca(trocaID, validar)
	if err != nil {
		return nil, err
	}

	if validar {
		transacao := &domain.BoicoinsTransacoes{
			ID:            uuid.New(),
			UsuarioID:     t.UsuarioID,
			Quantidade:    t.BoicoinsRecebidos,
			TipoTransacao: "recebimento_troca",
			Descricao:     "Recebimento de Boicoins por troca de item",
			DataTransacao: time.Now(),
			TrocaID:       &t.ID,
		}

		if _, err := tuc.boicoinRepo.Create(transacao); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (tuc *TrocaUseCase) notificarAdministrador(troca *domain.Troca) error {
	m, err := prepararEmailNotificacao(troca)
	if err != nil {
		return fmt.Errorf("erro ao preparar o e-mail: %w", err)
	}

	d := mail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		logrus.WithError(err).Error("erro ao enviar e-mail")
		return fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	logrus.Info("E-mail de notificação enviado com sucesso!")
	return nil
}

func prepararEmailNotificacao(troca *domain.Troca) (*mail.Message, error) {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", "logancardoso4@gmail.com")
	m.SetHeader("Subject", "Nova doação pendente - ID: "+troca.ID.String())

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Nova Doação Pendente</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333;
					line-height: 1.6;
				}
				.container {
					width: 80%%;
					max-width: 600px;
					margin: 0 auto;
					background: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				h2 {
					color: #2c3e50;
				}
				.details {
					background-color: #ecf0f1;
					padding: 10px;
					border-radius: 5px;
					margin-top: 15px;
				}
				.details p {
					margin: 5px 0;
				}
				.button {
					display: inline-block;
					background-color: #3498db;
					color: #ffffff;
					padding: 10px 20px;
					text-align: center;
					text-decoration: none;
					border-radius: 5px;
					margin-top: 20px;
					font-size: 16px;
				}
				.button:hover {
					background-color: #2980b9;
				}
				.footer {
					margin-top: 20px;
					text-align: center;
					color: #777;
					font-size: 12px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Nova doação pendente de validação</h2>
				<p>Uma nova doação foi registrada e está aguardando validação. Veja os detalhes abaixo:</p>
				
				<div class="details">
					<p><strong>ID da Doação:</strong> %s</p>
					<p><strong>ID do Usuário:</strong> %s</p>
					<p><strong>ID do Item:</strong> %s</p>
					<p><strong>Quantidade:</strong> %d</p>
					<p><strong>Data:</strong> %s</p>
				</div>
				
				<a href="https://seuapp.com/admin/validar-troca/%s" class="button">Validar Doação</a>

				<div class="footer">
					<p>Este é um e-mail automático. Por favor, não responda.</p>
				</div>
			</div>
		</body>
		</html>
	`,
		troca.ID.String(),
		troca.UsuarioID.String(),
		troca.ItemTrocaID.String(),
		troca.Quantidade,
		troca.DataTroca.Format(time.RFC1123),
		troca.ID.String(),
	)

	m.SetBody("text/html", body)
	return m, nil
}

func (tuc *TrocaUseCase) GetTroca(idTroca uuid.UUID) (*domain.Troca, error) {
	troca, err := tuc.trocaRepo.GetByID(idTroca)
	if err != nil {
		return nil, err
	}

	return troca, nil
}

func (tuc *TrocaUseCase) TodosItensTroca() ([]*domain.ItemTroca, error) {
	itensTroca, err := tuc.ItemTrocaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return itensTroca, nil
}

func (tuc *TrocaUseCase) CriarItemTroca(ItemTrocaRequest *domain.ItemTroca) (*domain.ItemTroca, error) {
	ItemTroca := &domain.ItemTroca{
		ID:                 uuid.New(),
		NomeItem:           ItemTrocaRequest.NomeItem,
		Descricao:          ItemTrocaRequest.Descricao,
		UnidadeMedida:      ItemTrocaRequest.UnidadeMedida,
		BoicoinsPorUnidade: ItemTrocaRequest.BoicoinsPorUnidade,
	}

	return tuc.ItemTrocaRepo.Create(ItemTroca)
}

func (tuc *TrocaUseCase) AtualizaItemTroca(ItemTrocaRequest *domain.ItemTroca) (*domain.ItemTroca, error) {
	ItemTroca, err := tuc.ItemTrocaRepo.GetByID(ItemTrocaRequest.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item de doação não encontrado")
		}
		return nil, err
	}

	if ItemTrocaRequest.Descricao != "" {
		ItemTroca.Descricao = ItemTrocaRequest.Descricao
	}
	if ItemTrocaRequest.UnidadeMedida != "" {
		ItemTroca.UnidadeMedida = ItemTrocaRequest.UnidadeMedida
	}
	if ItemTrocaRequest.BoicoinsPorUnidade != 0 {
		ItemTroca.BoicoinsPorUnidade = ItemTrocaRequest.BoicoinsPorUnidade
	}

	return tuc.ItemTrocaRepo.Update(ItemTroca)
}

func (tuc *TrocaUseCase) DeletarItemTroca(id uuid.UUID) error {
	if err := tuc.ItemTrocaRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item de doação não encontrado")
		}
		return err
	}
	return nil
}

func calculaBoicoins(valorUnidade float64, quantidade float64) (float64, error) {
	if valorUnidade < 0 || quantidade < 0 {
		return 0, errors.New("valorUnidade e quantidade devem ser maiores que zero")
	}
	return valorUnidade * float64(quantidade), nil
}

func (tuc *TrocaUseCase) IniciarCronJobExpiracaoTroca() {
	c := cron.New()

	// Executar a cada dia às 2 da manhã para verificar e expirar trocas
	_, err := c.AddFunc("0 2 * * *", func() {
		logrus.Info("Cron job para expirar trocas iniciado.")
		if err := tuc.marcarTrocasExpiradas(); err != nil {
			logrus.Error("Erro ao deletar trocas expiradas: ", err)
		} else {
			logrus.Info("Cron job para expirar trocas finalizado com sucesso.")
		}
	})
	if err != nil {
		logrus.Fatal("Erro ao configurar o cron job: ", err)
	}

	// Iniciar o cron job e registrar que está ativo
	c.Start()
	logrus.Info("Cron job para expirar trocas configurado e ativo.")
}

func (tuc *TrocaUseCase) marcarTrocasExpiradas() error {
	dataExpiracao := time.Now().AddDate(0, 0, -16)

	logrus.Infof("Buscando trocas criadas antes de %s para expiração.", dataExpiracao.Format("2006-01-02 15:04:05"))

	err := tuc.trocaRepo.DeletarTrocasCriadasAntesDe(dataExpiracao)
	if err != nil {
		logrus.Error("Erro ao deletar trocas expiradas: ", err)
		return err
	}

	logrus.Info("Processo de expiração completado. Todas as trocas expiradas foram deletadas.")
	return nil
}

//
//func (tuc *TrocaUseCase) IniciarCronJobExpiracaoTroca() {
//	c := cron.New()
//
//	// Para fins de teste: Executar a cada minuto para verificar se o cron job está funcionando corretamente
//	_, err := c.AddFunc("@every 1m", func() {
//		logrus.Info("Cron job para expirar trocas iniciado (teste com 1 minuto).")
//		if err := tuc.marcarTrocasExpiradas(); err != nil {
//			logrus.Error("Erro ao deletar trocas expiradas: ", err)
//		} else {
//			logrus.Info("Cron job para expirar trocas finalizado com sucesso (teste com 1 minuto).")
//		}
//	})
//	if err != nil {
//		logrus.Fatal("Erro ao configurar o cron job para teste: ", err)
//	}
//
//	// Iniciar o cron job e registrar que está ativo
//	c.Start()
//	logrus.Info("Cron job para expirar trocas configurado e ativo (teste com 1 minuto).")
//}
//
//func (tuc *TrocaUseCase) marcarTrocasExpiradas() error {
//	// Para fins de teste: Expirar trocas criadas 1 minuto antes
//	dataExpiracao := time.Now().Add(-1 * time.Minute)
//
//	logrus.Infof("Buscando trocas criadas antes de %s para expiração (teste com 1 minuto).", dataExpiracao.Format("2006-01-02 15:04:05"))
//
//	// Deletar todas as trocas criadas antes da data de expiração
//	err := tuc.trocaRepo.DeletarTrocasCriadasAntesDe(dataExpiracao)
//	if err != nil {
//		logrus.Error("Erro ao deletar trocas expiradas (teste): ", err)
//		return err
//	}
//
//	logrus.Info("Processo de expiração completado. Todas as trocas expiradas foram deletadas (teste).")
//	return nil
//}
