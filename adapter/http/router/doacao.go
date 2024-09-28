package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupDoacaoRoutes(router *gin.Engine, doacaoHandler *handler.DoacaoHandler) {
	protectedGroup := router.Group("/api/v1/doacoes")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("", doacaoHandler.AdicionaDoacao)
	}
}
