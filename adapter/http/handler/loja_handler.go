package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	minioClient "boi-marronzinho-api/minio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	idStr := c.Param("usuarioID")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário é inválido"})
		return
	}

	err = lh.LojaUseCase.AdicionarOuIncrementarItemCarrinho(usuarioID, itemCarrinho.ProdutoID, itemCarrinho.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item adicionado ou atualizado no carrinho com sucesso"})
}

func (lh *LojaHandler) ListarItensCarrinho(c *gin.Context) {
	idStr := c.Param("usuarioID")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário é inválido"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário é inválido"})
		return
	}
	err = lh.LojaUseCase.RemoverItemCarrinho(usuarioID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removido do carrinho com sucesso"})
}

func (lh *LojaHandler) AtualizarItemCarrinho(c *gin.Context) {
	var itemCarrinho struct {
		ProdutoID  uuid.UUID `json:"produtoId" validate:"required"`
		Quantidade int       `json:"quantidade" validate:"required,gte=1"`
	}
	if err := c.ShouldBindJSON(&itemCarrinho); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("usuarioID")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário é inválido"})
		return
	}

	err = lh.LojaUseCase.AtualizarQuantidadeItemCarrinho(usuarioID, itemCarrinho.ProdutoID, itemCarrinho.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quantidade do item atualizada com sucesso"})
}

func (lh *LojaHandler) FinalizarCompra(c *gin.Context) {
	idStr := c.Param("usuarioID")
	usuarioID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário é inválido"})
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
	var produtoDTO domain.Produto

	jsonData := c.PostForm("request")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'request' com os dados do evento é obrigatório."})
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &produtoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao parsear os dados JSON do evento."})
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

	imageFileName := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)

	imageUrl, err := minioClient.UploadFile(fileContent, imageFileName, "produtos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer upload da imagem"})
		return
	}

	produtoDTO.ImagemURL = imageUrl
	produtoDTO.CriadoEm = time.Now()
	produtoDTO.ID = uuid.New()

	produto, err := lh.LojaUseCase.AdicionarProduto(&produtoDTO)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto é inválido"})
		return
	}

	var updateData domain.Produto
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	produtoAtualizado, err := lh.LojaUseCase.EditaProduto(id, &updateData)
	if err != nil {
		if err.Error() == "produto não encontrado" {
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
