package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupTrocaRoutes(router *gin.Engine, pedidoHandler *handler.TrocaHandler) {
	protectedGroup := router.Group("/api/v1/trocas")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("", pedidoHandler.RealizarTroca)
		protectedGroup.POST("/scanner-troca/:id", auth.RoleAuthMiddleware(auth.GetRole()), pedidoHandler.GetTroca)
		protectedGroup.GET("/itens-troca", pedidoHandler.ItensDoacao)
		protectedGroup.POST("/adiciona-item-troca", auth.RoleAuthMiddleware(auth.GetRole()), pedidoHandler.CriaItemTroca)
		protectedGroup.PUT("/:id", auth.RoleAuthMiddleware(auth.GetRole()), pedidoHandler.AtualizaItemTroca)
		protectedGroup.DELETE("/:id", auth.RoleAuthMiddleware(auth.GetRole()), pedidoHandler.DeletaItemTroca)
		protectedGroup.POST("valida-troca", auth.RoleAuthMiddleware(auth.GetRole()), pedidoHandler.ValidaTroca)
	}
}
