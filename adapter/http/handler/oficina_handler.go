package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"boi-marronzinho-api/dto"
	"net/http"

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
	var oficinaDTO *domain.Oficinas
	if err := c.ShouldBindJSON(&oficinaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oficinaDTO.ID = uuid.New()

	oficina, err := oh.OficinaUseCase.CriaOficina(oficinaDTO)
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

	c.JSON(http.StatusOK, oficinas)
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

	if err := c.ShouldBindJSON(&voucherDTO); err!= nil {
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