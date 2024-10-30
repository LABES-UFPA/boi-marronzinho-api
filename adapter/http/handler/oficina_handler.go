package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	minioClient "boi-marronzinho-api/minio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	minio "github.com/minio/minio-go/v7"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OficinaHandler struct {
	OficinaUseCase *usecase.OficinaUseCase
	MinioClient    *minio.Client
}

func NewOficinaHandler(ouc *usecase.OficinaUseCase, minioClient *minio.Client) *OficinaHandler {
	return &OficinaHandler{OficinaUseCase: ouc, MinioClient: minioClient}
}

func (oh *OficinaHandler) CriaOficina(c *gin.Context) {
	var oficinaDTO domain.Oficinas

	jsonData := c.PostForm("request")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'request' com os dados da oficina é obrigatório."})
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &oficinaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao parsear os dados JSON da oficina."})
		return
	}

	dataEvento, err := time.Parse(time.RFC3339, oficinaDTO.DataEvento.Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use o formato RFC3339."})
		return
	}
	oficinaDTO.DataEvento = dataEvento

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao receber o arquivo de imagem"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo de imagem"})
		return
	}
	defer fileContent.Close()

	imageFileName := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)

	imageUrl, err := minioClient.UploadFile(fileContent, imageFileName, "oficinas")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer upload da imagem"})
		return
	}

	oficinaDTO.ImagemUrl = imageUrl
	oficinaDTO.ID = uuid.New()
	oficinaDTO.CriadoEm = time.Now()

	oficina, err := oh.OficinaUseCase.CriaOficina(&oficinaDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, oficina)
}

func (oh *OficinaHandler) ListaOficinas(c *gin.Context) {
	oficinas, err := oh.OficinaUseCase.ListaOficinas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var oficinasResponse []dto.OficinaResponse
	for _, oficina := range oficinas {
		oficinaResponse := dto.OficinaResponse{
			ID:                  oficina.ID,
			Nome:                oficina.Nome,
			Descricao:           oficina.Descricao,
			PrecoBoicoins:       oficina.PrecoBoicoins,
			PrecoReal:           oficina.PrecoReal,
			DataEvento:          oficina.DataEvento,
			LimiteParticipantes: oficina.LimiteParticipantes,
			ParticipantesAtual:  oficina.ParticipantesAtual,
			LinkEndereco:        oficina.LinkEndereco,
			Imagem:              oficina.ImagemUrl,
		}
		oficinasResponse = append(oficinasResponse, oficinaResponse)
	}

	c.JSON(http.StatusOK, oficinasResponse)
}

func (oh *OficinaHandler) UpdateOficina(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da oficina é inválido"})
		return
	}

	jsonData := c.PostForm("request")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'request' com os dados da oficina é obrigatório."})
		return
	}

	var updateData domain.Oficinas
	if err := json.Unmarshal([]byte(jsonData), &updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao parsear os dados JSON da oficina."})
		return
	}

	dataEvento, err := time.Parse(time.RFC3339, updateData.DataEvento.Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use o formato RFC3339."})
		return
	}
	updateData.DataEvento = dataEvento

	file, err := c.FormFile("file")
	if err == nil {
		updatedOficina, err := oh.OficinaUseCase.UpdateOficinaWithFile(id, &updateData, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedOficina)
		return
	}

	updatedOficina, err := oh.OficinaUseCase.UpdateOficina(id, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedOficina)
}

func (oh *OficinaHandler) DeleteOficina(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da oficina é inválido"})
		return
	}

	err = oh.OficinaUseCase.DeleteOficina(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "oficina deletada com sucesso!"})
}

func (oh *OficinaHandler) InscricaoOficina(c *gin.Context) {
	var inscricaoDTO *domain.ParticipanteOficina
	if err := c.ShouldBindJSON(&inscricaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagamentoEmBoicoins := c.DefaultQuery("boicoins", "true") == "true"

	inscricao, err := oh.OficinaUseCase.InscricaoOficina(inscricaoDTO, pagamentoEmBoicoins)
	if err != nil {
		if err.Error() == "saldo de Boicoins insuficiente" || err.Error() == "não há mais vagas disponíveis para esta oficina" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "usuário não encontrado" || err.Error() == "oficina não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, inscricao)
}

func (oh *OficinaHandler) ListaTicketsPorUsuario(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	tickets, err := oh.OficinaUseCase.ListarTicketsPorUsuario(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (oh *OficinaHandler) ScannerQRCode(c *gin.Context) {
	var voucherDTO *dto.ValidaVoucherDTO

	if err := c.ShouldBindJSON(&voucherDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valido, err := oh.OficinaUseCase.ValidaVoucher(&voucherDTO.IdVoucher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, valido)
}
