package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

type CarteiraHandler struct {
	CarteiraUseCase *usecase.CarteiraUseCase
}

func NewCarteiraHandler(cuc *usecase.CarteiraUseCase) *CarteiraHandler {
	return &CarteiraHandler{CarteiraUseCase: cuc}
}

func (ch *CarteiraHandler) AdicionaTransacao(c *gin.Context) {
	var request struct {
		UsuarioID       string `json:"usuarioID"`
		Quantidade      float64 `json:"quantidade"`
		Descricao       string `json:"descricao"`
		DoacaoID        string `json:"doacaoID"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuarioUUID, err := uuid.Parse(request.UsuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	var doacaoUUID *uuid.UUID
	if request.DoacaoID != "" {
		parsedID, err := uuid.Parse(request.DoacaoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de doação inválido"})
			return
		}
		doacaoUUID = &parsedID
	}

	transacao := &domain.BoicoinsTransacoes{
		ID:            uuid.New(),
		UsuarioID:     usuarioUUID,
		Quantidade:    request.Quantidade,
		TipoTransacao: "recebimento_doacao",
		Descricao:     request.Descricao,
		DataTransacao: time.Now(),
		DoacaoID:      doacaoUUID,
	}

	

	if err := ch.CarteiraUseCase.CriaTransacao(transacao); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transação adicionada com sucesso"})
}
