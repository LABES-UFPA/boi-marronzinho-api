package handler

import (
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	UserUseCase *usecase.UsuarioUseCase
}

func NewUserHandler(uuc *usecase.UsuarioUseCase) *UsuarioHandler {
	return &UsuarioHandler{UserUseCase: uuc}
}

func (uh *UsuarioHandler) CreateUser(c *gin.Context) {
	var userDTO *domain.Usuario
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserUseCase.CreateUser(userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
