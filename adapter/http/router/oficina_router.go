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
		protectedGroup.GET("/meus-tickets/:id", oficinaHandler.ListaTicketsPorUsuario)

		protectedGroup.POST("/cria-oficinas", auth.RoleAuthMiddleware(auth.GetRole()), oficinaHandler.CriaOficina)
		protectedGroup.POST("/scanner-voucher", auth.RoleAuthMiddleware(auth.GetRole()), oficinaHandler.ScannerQRCode)
		protectedGroup.DELETE("/deleta-oficinas/:id", auth.RoleAuthMiddleware(auth.GetRole()), oficinaHandler.DeleteOficina)
		protectedGroup.PUT("/atualiza-oficinas/:id", auth.RoleAuthMiddleware(auth.GetRole()), oficinaHandler.UpdateOficina)
	}
}
