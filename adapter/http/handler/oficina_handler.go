package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OficinaHandler struct {
	OficinaUseCase *usecase.OficinaUseCase
}

func NewOficinaHandler(ouc *usecase.OficinaUseCase) *OficinaHandler {
	return &OficinaHandler{OficinaUseCase: ouc}
}

func (oh *OficinaHandler) CriaOficina(c *gin.Context) {
	var oficinaDTO domain.Oficinas

	if err := c.ShouldBind(&oficinaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar dados da oficina: " + err.Error()})
		return
	}

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

	imageData, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler o arquivo de imagem"})
		return
	}

	oficinaDTO.Imagem = imageData

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
			Imagem:              base64.StdEncoding.EncodeToString(oficina.Imagem),
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

	var updateData domain.Oficinas
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	updatedUser, err := oh.OficinaUseCase.UpdateOficina(id, &updateData)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (oh *OficinaHandler) DeleteOficina(c *gin.Context) {
	var oficinaDTO *dto.Teste
	if err := c.ShouldBindJSON(&oficinaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := oh.OficinaUseCase.DeleteUser(uuid.MustParse(oficinaDTO.ID))
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
