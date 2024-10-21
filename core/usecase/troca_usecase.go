package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TrocaUseCase struct {
	doacaoRepo    repository.Repository[domain.Troca]
	ItemTrocaRepo repository.Repository[domain.ItemTroca]
	userRepo      repository.UserRepository
	boicoinRepo   repository.BoicoinRepository
}

func NewTrocaUseCase(
	trocaRepo repository.Repository[domain.Troca],
	ItemTrocaRepo repository.Repository[domain.ItemTroca],
	userRepo repository.UserRepository,
	boicoinRepo repository.BoicoinRepository,
) *TrocaUseCase {
	return &TrocaUseCase{
		doacaoRepo:    trocaRepo,
		ItemTrocaRepo: ItemTrocaRepo,
		userRepo:      userRepo,
		boicoinRepo:   boicoinRepo,
	}
}

func (duc *TrocaUseCase) RealizarTroca(trocaRequest *domain.Troca) (*domain.Troca, error) {
	ItemTroca, err := duc.ItemTrocaRepo.GetByID(trocaRequest.ItemTrocaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item de doação não encontrado")
		}
		return nil, err
	}

	boicoinsRecebidos, err := calculaBoicoins(ItemTroca.BoicoinsPorUnidade, trocaRequest.Quantidade)
	if err != nil {
		return nil, err
	}

	doacao := &domain.Troca{
		ID:                uuid.New(),
		UsuarioID:         trocaRequest.UsuarioID,
		ItemTrocaID:       trocaRequest.ItemTrocaID,
		Quantidade:        trocaRequest.Quantidade,
		BoicoinsRecebidos: boicoinsRecebidos,
		DataDoacao:        time.Now(),
		Status:            "pendente",
	}

	createdDoacao, err := duc.doacaoRepo.Create(doacao)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := duc.notificarAdministrador(createdDoacao); err != nil {
			logrus.Info("Erro ao notificar administrador" + err.Error())
		}
	}()

	return createdDoacao, nil
}

func (duc *TrocaUseCase) ValidaTroca(doacaoID string, validar bool) (*domain.Troca, error) {
	doacao, err := duc.doacaoRepo.GetByID(uuid.MustParse(doacaoID))
	if err != nil {
		return nil, err
	}

	if validar {
		doacao.Status = "validada"

		transacao := &domain.BoicoinsTransacoes{
			ID:            uuid.New(),
			UsuarioID:     doacao.UsuarioID,
			Quantidade:    +float64(doacao.BoicoinsRecebidos),
			TipoTransacao: "recebimento_doacao",
			Descricao:     "Recebimento de Boicoins por doação de item",
			DataTransacao: time.Now(),
			DoacaoID:      &doacao.ID,
		}

		if _, err := duc.boicoinRepo.Create(transacao); err != nil {
			return nil, err
		}
	} else {
		doacao.Status = "rejeitada"
	}
	doacaoAtualizada, err := duc.doacaoRepo.Update(doacao)
	if err != nil {
		return nil, err
	}

	return doacaoAtualizada, nil
}

func (duc *TrocaUseCase) notificarAdministrador(doacao *domain.Troca) error {
	m, err := prepararEmailNotificacao(doacao)
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

func prepararEmailNotificacao(doacao *domain.Troca) (*mail.Message, error) {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", "logancardoso4@gmail.com")
	m.SetHeader("Subject", "Nova doação pendente - ID: "+doacao.ID.String())

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
				
				<a href="https://seuapp.com/admin/validar-doacao/%s" class="button">Validar Doação</a>

				<div class="footer">
					<p>Este é um e-mail automático. Por favor, não responda.</p>
				</div>
			</div>
		</body>
		</html>
	`,
		doacao.ID.String(),
		doacao.UsuarioID.String(),
		doacao.ItemTrocaID.String(),
		doacao.Quantidade,
		doacao.DataDoacao.Format(time.RFC1123),
		doacao.ID.String(),
	)

	m.SetBody("text/html", body)
	return m, nil
}

func (duc *TrocaUseCase) TodosItensTroca() ([]*domain.ItemTroca, error) {
	itensTroca, err := duc.ItemTrocaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return itensTroca, nil
}

func (duc *TrocaUseCase) CriarItemTroca(ItemTrocaRequest *domain.ItemTroca) (*domain.ItemTroca, error) {
	ItemTroca := &domain.ItemTroca{
		ID:                 uuid.New(),
		Descricao:          ItemTrocaRequest.Descricao,
		UnidadeMedida:      ItemTrocaRequest.UnidadeMedida,
		BoicoinsPorUnidade: ItemTrocaRequest.BoicoinsPorUnidade,
	}

	return duc.ItemTrocaRepo.Create(ItemTroca)
}

func (duc *TrocaUseCase) AtualizaItemTroca(ItemTrocaRequest *domain.ItemTroca) (*domain.ItemTroca, error) {
	ItemTroca, err := duc.ItemTrocaRepo.GetByID(ItemTrocaRequest.ID)
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

	return duc.ItemTrocaRepo.Update(ItemTroca)
}

func (duc *TrocaUseCase) DeletarItemTroca(id uuid.UUID) error {
	if err := duc.ItemTrocaRepo.Delete(id); err != nil {
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
