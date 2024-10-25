package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupCarteiraRoutes(router *gin.Engine, carteiraHandler *handler.CarteiraHandler) {
	protectedGroup := router.Group("/api/v1/carteiras")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("/adicionar-transacao", auth.RoleAuthMiddleware(auth.GetRole()), carteiraHandler.AdicionaTransacao)
	}
}