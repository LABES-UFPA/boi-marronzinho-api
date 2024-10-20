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

type DoacaoUseCase struct {
	doacaoRepo     repository.Repository[domain.Doacoes]
	itemDoacaoRepo repository.Repository[domain.ItemDoacao]
	userRepo       repository.UserRepository
	boicoinRepo    repository.BoicoinRepository
}

func NewDoacaoUseCase(
	doacaoRepo repository.Repository[domain.Doacoes],
	itemDoacaoRepo repository.Repository[domain.ItemDoacao],
	userRepo repository.UserRepository,
	boicoinRepo repository.BoicoinRepository,
) *DoacaoUseCase {
	return &DoacaoUseCase{
		doacaoRepo:     doacaoRepo,
		itemDoacaoRepo: itemDoacaoRepo,
		userRepo:       userRepo,
		boicoinRepo:    boicoinRepo,
	}
}

func (duc *DoacaoUseCase) AdicionaDoacao(doacaoRequest *domain.Doacoes) (*domain.Doacoes, error) {
	itemDoacao, err := duc.itemDoacaoRepo.GetByID(doacaoRequest.ItemDoacaoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item de doação não encontrado")
		}
		return nil, err
	}

	boicoinsRecebidos, err := calculaBoicoins(itemDoacao.BoicoinsPorUnidade, doacaoRequest.Quantidade)
	if err != nil {
		return nil, err
	}

	doacao := &domain.Doacoes{
		ID:                uuid.New(),
		UsuarioID:         doacaoRequest.UsuarioID,
		ItemDoacaoID:      doacaoRequest.ItemDoacaoID,
		Quantidade:        doacaoRequest.Quantidade,
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

func (duc *DoacaoUseCase) ValidaDoacao(doacaoID string, validar bool) (*domain.Doacoes, error) {
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

// Função para notificar o administrador por e-mail
func (duc *DoacaoUseCase) notificarAdministrador(doacao *domain.Doacoes) error {
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

func prepararEmailNotificacao(doacao *domain.Doacoes) (*mail.Message, error) {
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
		doacao.ItemDoacaoID.String(),
		doacao.Quantidade,
		doacao.DataDoacao.Format(time.RFC1123),
		doacao.ID.String(),
	)

	m.SetBody("text/html", body)
	return m, nil
}


func (duc *DoacaoUseCase) TodosItensDoacao() ([]*domain.ItemDoacao, error) {
	itensDoacao, err := duc.itemDoacaoRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return itensDoacao, nil
}

func (duc *DoacaoUseCase) CriarItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	itemDoacao := &domain.ItemDoacao{
		ID:                 uuid.New(),
		Descricao:          itemDoacaoRequest.Descricao,
		UnidadeMedida:      itemDoacaoRequest.UnidadeMedida,
		BoicoinsPorUnidade: itemDoacaoRequest.BoicoinsPorUnidade,
	}

	return duc.itemDoacaoRepo.Create(itemDoacao)
}

func (duc *DoacaoUseCase) AtualizaItemDoacao(itemDoacaoRequest *domain.ItemDoacao) (*domain.ItemDoacao, error) {
	itemDoacao, err := duc.itemDoacaoRepo.GetByID(itemDoacaoRequest.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item de doação não encontrado")
		}
		return nil, err
	}

	if itemDoacaoRequest.Descricao != "" {
		itemDoacao.Descricao = itemDoacaoRequest.Descricao
	}
	if itemDoacaoRequest.UnidadeMedida != "" {
		itemDoacao.UnidadeMedida = itemDoacaoRequest.UnidadeMedida
	}
	if itemDoacaoRequest.BoicoinsPorUnidade != 0 {
		itemDoacao.BoicoinsPorUnidade = itemDoacaoRequest.BoicoinsPorUnidade
	}

	return duc.itemDoacaoRepo.Update(itemDoacao)
}

func (duc *DoacaoUseCase) DeletarItemDoacao(id uuid.UUID) error {
	if err := duc.itemDoacaoRepo.Delete(id); err != nil {
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
