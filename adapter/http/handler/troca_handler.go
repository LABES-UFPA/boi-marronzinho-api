package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrocaHandler struct {
	TrocaUseCase *usecase.TrocaUseCase
}

func NewTrocaHandler(duc *usecase.TrocaUseCase) *TrocaHandler {
	return &TrocaHandler{TrocaUseCase: duc}
}

func (th *TrocaHandler) RealizarTroca(c *gin.Context) {
	var trocaDTO *domain.Troca
	if err := c.ShouldBindJSON(&trocaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	troca, qrcodeValidacao, err := th.TrocaUseCase.RealizarTroca(trocaDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"troca":        troca,
		"qrcode_base64": qrcodeValidacao,
	})
}

func (th *TrocaHandler) GetTroca(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da troca é inválido"})
		return
	}
	troca, err := th.TrocaUseCase.GetTroca(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, troca)
}

func (th *TrocaHandler) ValidaTroca(c *gin.Context) {
	var request struct {
		TrocaID string `json:"trocaID"`
		Validar bool   `json:"validar"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valido, err := th.TrocaUseCase.ValidaTroca(uuid.MustParse(request.TrocaID), request.Validar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, valido)
}

func (th *TrocaHandler) CriaItemTroca(c *gin.Context) {
	var itemTrocaDTO *domain.ItemTroca
	if err := c.ShouldBindJSON(&itemTrocaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemTroca, err := th.TrocaUseCase.CriarItemTroca(itemTrocaDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemTroca)
}

func (th *TrocaHandler) ItensDoacao(c *gin.Context) {
	itensTroca, err := th.TrocaUseCase.TodosItensTroca()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itensTroca)
}

func (th *TrocaHandler) AtualizaItemTroca(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de troca é inválido"})
		return
	}

	var itemTrocaDTO *domain.ItemTroca
	if err := c.ShouldBindJSON(&itemTrocaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemTrocaDTO.ID = id
	itemTroca, err := th.TrocaUseCase.AtualizaItemTroca(itemTrocaDTO)
	if err != nil {
		if err.Error() == "item de troca não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "item de troca não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, itemTroca)
}

func (th *TrocaHandler) DeletaItemTroca(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de doação é inválido"})
		return
	}

	if err := th.TrocaUseCase.DeletarItemTroca(id); err != nil {
		if err.Error() == "item de troca não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "item de troca não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item de troca deletado com sucesso"})
}