package usecase

import (
	"boi-marronzinho-api/adapter/repository"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	"boi-marronzinho-api/minio"
	minioClient "boi-marronzinho-api/minio"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OficinaUseCase struct {
	oficinaRepo repository.OficinaRepository
	usuarioRepo repository.UserRepository
}

func NewOficinaUseCase(oficinaRepo repository.OficinaRepository, usuarioRepo repository.UserRepository) *OficinaUseCase {
	return &OficinaUseCase{
		oficinaRepo: oficinaRepo,
		usuarioRepo: usuarioRepo,
	}
}

func (o *OficinaUseCase) CriaOficina(oficinaRequest *domain.Oficinas) (*domain.Oficinas, error) {
	// if err := oficinaRequest.Validate(); err != nil {
	//     return nil, err
	// }
	oficina, err := o.oficinaRepo.Create(oficinaRequest)
	if err != nil {
		return nil, err
	}

	return oficina, nil
}

func (o *OficinaUseCase) ListaOficinas() ([]*domain.Oficinas, error) {
	return o.oficinaRepo.GetAll()
}

func (o *OficinaUseCase) UpdateOficina(id uuid.UUID, updateData *domain.Oficinas) (*domain.Oficinas, error) {
	oficina, err := o.oficinaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if updateData.Nome != "" {
		oficina.Nome = updateData.Nome
	}
	if updateData.LimiteParticipantes != nil {
		oficina.LimiteParticipantes = updateData.LimiteParticipantes
	}

	if _, err = o.oficinaRepo.Update(oficina); err != nil {
		return nil, errors.New("falha ao atualizar o usuário")
	}

	return oficina, nil
}

func (o *OficinaUseCase) UpdateOficinaWithFile(id uuid.UUID, updateData *domain.Oficinas, file *multipart.FileHeader) (*domain.Oficinas, error) {
    oficina, err := o.oficinaRepo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("oficina não encontrada")
        }
        return nil, err
    }

    if updateData.Nome != "" {
        oficina.Nome = updateData.Nome
    }
    if updateData.LimiteParticipantes != nil {
        oficina.LimiteParticipantes = updateData.LimiteParticipantes
    }
    oficina.DataEvento = updateData.DataEvento

    if file != nil {
        oldFileName := oficina.ImagemUrl

        src, err := file.Open()
        if err != nil {
            return nil, errors.New("falha ao abrir o arquivo de imagem")
        }
        defer src.Close()

        newFileName := uuid.New().String()

        imageURL, err := minio.UploadFile(src, newFileName, "oficinas")
        if err != nil {
            return nil, errors.New("falha ao fazer upload da nova imagem")
        }

        oficina.ImagemUrl = imageURL

        if oldFileName != "" {
            if err := minio.DeleteFile(oldFileName, "oficinas"); err != nil {
                return nil, fmt.Errorf("falha ao deletar a imagem antiga: %v", err)
            }
        }
    }

    updatedOficina, err := o.oficinaRepo.Update(oficina)
    if err != nil {
        return nil, errors.New("falha ao atualizar a oficina")
    }

    return updatedOficina, nil
}

func (o *OficinaUseCase) DeleteOficina(id uuid.UUID) error {
	oficina, err := o.oficinaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("oficina não encontrada")
		}
		return errors.New("erro ao buscar oficina: " + err.Error())
	}

	if err := o.oficinaRepo.DeleteOficina(id); err != nil {
		return errors.New("erro ao deletar oficina: " + err.Error())
	}

	logrus.Infof("Oficina ID %s deletada com sucesso no caso de uso", id)

	fileName := oficina.ImagemUrl
	bucketName := "oficinas"

	if err := minioClient.DeleteFile(fileName, bucketName); err != nil {
		logrus.Warnf("Falha ao deletar imagem associada à oficina ID %s: %v", id, err)
	} else {
		logrus.Infof("Imagem %s deletada com sucesso do bucket %s", fileName, bucketName)
	}

	return nil
}

func (o *OficinaUseCase) InscricaoOficina(inscricaoRequest *domain.ParticipanteOficina, pagamentoEmBoicoins bool) (*domain.ParticipanteOficina, error) {
	oficina, err := o.oficinaRepo.GetByID(inscricaoRequest.OficinaID)
	if err != nil {
		return nil, errors.New("oficina não encontrada")
	}

	if oficina.ParticipantesAtual >= *oficina.LimiteParticipantes {
		return nil, errors.New("não há mais vagas disponíveis para esta oficina")
	}

	usuario, err := o.usuarioRepo.GetByID(inscricaoRequest.UsuarioID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if pagamentoEmBoicoins {
		if usuario.SaldoBoicoins < float32(oficina.PrecoBoicoins) {
			return nil, errors.New("saldo de Boicoins insuficiente")
		}
	}

	po, err := o.oficinaRepo.InscreverParticipante(inscricaoRequest.OficinaID, usuario, pagamentoEmBoicoins)
	if err != nil {
		return nil, err
	}

	return po, nil
}

func (o *OficinaUseCase) ListarTicketsPorUsuario(usuarioID uuid.UUID) ([]dto.VoucherResponseDTO, error) {
	tickets, err := o.oficinaRepo.GetTicketsByUsuarioID(usuarioID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (o *OficinaUseCase) ValidaVoucher(codigoVoucher *string) (*dto.VoucherValidation, error) {

	validado, err := o.oficinaRepo.ValidaVoucher(codigoVoucher)
	if err != nil {
		return nil, err
	}

	validadoResponse := &dto.VoucherValidation{
		IDVoucher: *codigoVoucher,
		Validado:  validado.Validado,
	}

	return validadoResponse, nil
}
