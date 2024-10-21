package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupTrocaRoutes(router *gin.Engine, doacaoHandler *handler.TrocaHandler) {
	protectedGroup := router.Group("/api/v1/trocas")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("", doacaoHandler.RealizarTroca)
		protectedGroup.GET("/itens-troca", doacaoHandler.ItensDoacao)
		protectedGroup.POST("/adiciona-item-troca", auth.RoleAuthMiddleware(auth.GetRole()), doacaoHandler.CriaItemTroca)
		protectedGroup.PUT("/:id", auth.RoleAuthMiddleware(auth.GetRole()), doacaoHandler.AtualizaItemTroca)
		protectedGroup.DELETE("/:id", auth.RoleAuthMiddleware(auth.GetRole()), doacaoHandler.DeletaItemTroca)
		protectedGroup.POST("valida-troca", auth.RoleAuthMiddleware(auth.GetRole()), doacaoHandler.ValidaTroca)
	}
}
