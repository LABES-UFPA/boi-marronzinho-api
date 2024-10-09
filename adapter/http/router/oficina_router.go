package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupOficinaRoutes(router *gin.Engine, oficinaHandler *handler.OficinaHandler) {
	protectedGroup := router.Group("/api/v1/oficinas")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("/inscricao", oficinaHandler.InscricaoOficina)
		protectedGroup.GET("/lista-oficinas", oficinaHandler.ListaOficinas)
		protectedGroup.GET("/meus-tickets", oficinaHandler.ListaTicketsPorUsuario)

		protectedGroup.POST("/cria-oficinas", auth.RoleAuthMiddleware(getRole()), oficinaHandler.CriaOficina)
		protectedGroup.POST("/scanner-voucher", auth.RoleAuthMiddleware(getRole()), oficinaHandler.ScannerQRCode)
	}
}
