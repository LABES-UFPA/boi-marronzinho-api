package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LojaHandler struct {
	LojaUseCase *usecase.LojaUseCase
}

func NewLojaHandler(luc *usecase.LojaUseCase) *LojaHandler {
	return &LojaHandler{LojaUseCase: luc}
}

func (lh *LojaHandler) AdicionarItemCarrinho(c *gin.Context) {
	var itemCarrinho struct {
		ProdutoID  uuid.UUID `json:"produtoId" validate:"required"`
		Quantidade int       `json:"quantidade" validate:"required,gte=1"`
	}
	if err := c.ShouldBindJSON(&itemCarrinho); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuario é inválido"})
		return
	}

	err = lh.LojaUseCase.AdicionarItemCarrinho(usuarioID, itemCarrinho.ProdutoID, itemCarrinho.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item adicionado ao carrinho com sucesso"})
}

func (lh *LojaHandler) ListarItensCarrinho(c *gin.Context) {
	idStr := c.Param("id")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuario é inválido"})
		return
	}

	itens, err := lh.LojaUseCase.ListarItensCarrinho(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itens)
}

func (lh *LojaHandler) RemoverItemCarrinho(c *gin.Context) {
	itemIDStr := c.Param("itemID")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do item inválido"})
		return
	}

	idStr := c.Param("usuarioID")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuario é inválido"})
		return
	}
	err = lh.LojaUseCase.RemoverItemCarrinho(usuarioID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removido do carrinho com sucesso"})
}

// Finalizar Compra Handler

func (lh *LojaHandler) FinalizarCompra(c *gin.Context) {
	idStr := c.Param("id")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuario é inválido"})
		return
	}
	pedido, err := lh.LojaUseCase.FinalizarCompra(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

func (lh *LojaHandler) AdicionarProduto(c *gin.Context) {
	var produtoDTO *domain.Produto
	if err := c.ShouldBindJSON(&produtoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	produto, err := lh.LojaUseCase.AdicionarProduto(produtoDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, produto)
}

func (lh *LojaHandler) EditaProduto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuario é inválido"})
		return
	}

	var updateData domain.Produto
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	produtoAtualizado, err := lh.LojaUseCase.EditaProduto(id, &updateData)
	if err != nil {
		if err.Error() == "produto not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "produto não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, produtoAtualizado)
}

func (lh *LojaHandler) RemoveProduto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto é inválido"})
		return
	}

	err = lh.LojaUseCase.RemoveProduto(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "produto deletado com sucesso!"})
}

func (lh *LojaHandler) ListaProdutos(c *gin.Context) {
	produtos, err := lh.LojaUseCase.ListaProdutos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, produtos)
}
