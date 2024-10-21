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

func (dh *TrocaHandler) RealizarTroca(c *gin.Context) {
	var trocaDTO *domain.Troca
	if err := c.ShouldBindJSON(&trocaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	troca, err := dh.TrocaUseCase.RealizarTroca(trocaDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, troca)
}

func (dh *TrocaHandler) ValidaTroca(c *gin.Context) {
	var request struct {
		TrocaID string `json:"trocaID"`
		Validar bool   `json:"validar"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := dh.TrocaUseCase.ValidaTroca(request.TrocaID, request.Validar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doação processada com sucesso."})
}

func (dh *TrocaHandler) CriaItemTroca(c *gin.Context) {
	var itemTrocaDTO *domain.ItemTroca
	if err := c.ShouldBindJSON(&itemTrocaDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemTroca, err := dh.TrocaUseCase.CriarItemTroca(itemTrocaDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemTroca)
}

func (dh *TrocaHandler) ItensDoacao(c *gin.Context) {
	itensTroca, err := dh.TrocaUseCase.TodosItensTroca()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itensTroca)
}

func (dh *TrocaHandler) AtualizaItemTroca(c *gin.Context) {
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
	itemTroca, err := dh.TrocaUseCase.AtualizaItemTroca(itemTrocaDTO)
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

func (dh *TrocaHandler) DeletaItemTroca(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de doação é inválido"})
		return
	}

	if err := dh.TrocaUseCase.DeletarItemTroca(id); err != nil {
		if err.Error() == "item de troca não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "item de troca não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item de troca deletado com sucesso"})
}

// Captura todos os itens de doação
// func (dh *DoacaoHandler) CapturaTodosItensDoacao(c *gin.Context) {
//     itensDoacao, err := dh.TrocaUseCase.CapturaTodosItensDoacao()
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, itensDoacao)
// }

// Captura um item de doação específico pelo ID
// func (dh *DoacaoHandler) CapturaItemDoacao(c *gin.Context) {
//     idStr := c.Param("id")
//     id, err := uuid.Parse(idStr)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de doação é inválido"})
//         return
//     }

//     itemDoacao, err := dh.TrocaUseCase.CapturaItemDoacao(id)
//     if err != nil {
//         if err.Error() == "item de doação não encontrado" {
//             c.JSON(http.StatusNotFound, gin.H{"error": "item de doação não encontrado"})
//         } else {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         }
//         return
//     }

//     c.JSON(http.StatusOK, itemDoacao)
// }
