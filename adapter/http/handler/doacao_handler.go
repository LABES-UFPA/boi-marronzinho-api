package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (dh *DoacaoHandler) ValidaDoacao(c *gin.Context) {
	var request struct {
		DoacaoID string `json:"doacaoID"`
		Validar  bool   `json:"validar"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := dh.DoacaoUseCase.ValidaDoacao(request.DoacaoID, request.Validar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doação processada com sucesso."})
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

func (dh *DoacaoHandler) AtualizaItemDoacao(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de doação é inválido"})
		return
	}

	var itemDoacaoDTO *domain.ItemDoacao
	if err := c.ShouldBindJSON(&itemDoacaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemDoacaoDTO.ID = id // Garantir que o ID é o mesmo do parâmetro
	itemDoacao, err := dh.DoacaoUseCase.AtualizaItemDoacao(itemDoacaoDTO)
	if err != nil {
		if err.Error() == "item de doação não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "item de doação não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, itemDoacao)
}

func (dh *DoacaoHandler) DeletaItemDoacao(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item de doação é inválido"})
		return
	}

	if err := dh.DoacaoUseCase.DeletarItemDoacao(id); err != nil {
		if err.Error() == "item de doação não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "item de doação não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item de doação deletado com sucesso"})
}

// Captura todos os itens de doação
// func (dh *DoacaoHandler) CapturaTodosItensDoacao(c *gin.Context) {
//     itensDoacao, err := dh.DoacaoUseCase.CapturaTodosItensDoacao()
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

//     itemDoacao, err := dh.DoacaoUseCase.CapturaItemDoacao(id)
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
