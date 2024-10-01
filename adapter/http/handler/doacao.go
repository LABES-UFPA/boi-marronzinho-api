package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DoacaoHandler struct {
	DoacaoUseCase *usecase.DoacaoUseCase
}

func NewDoacaoHandler(duc *usecase.DoacaoUseCase) *DoacaoHandler {
	return &DoacaoHandler{DoacaoUseCase: duc}
}

func (dh *DoacaoHandler) AdicionaDoacao(c *gin.Context) {
	var doacaoDTO *domain.Doacoes
	if err := c.ShouldBindJSON(&doacaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doacao, err := dh.DoacaoUseCase.AdicionaDoacao(doacaoDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doacao)
}

func (dh *DoacaoHandler) CriaItemDoacao(c *gin.Context) {
	var itemDoacaoDTO *domain.ItemDoacao
	if err := c.ShouldBindJSON(&itemDoacaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemDoacao, err := dh.DoacaoUseCase.CriarItemDoacao(itemDoacaoDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemDoacao)
}
