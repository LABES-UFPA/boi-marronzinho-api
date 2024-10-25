package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"
		
	"github.com/gin-gonic/gin"
)

func SetupEventoRoutes(router *gin.Engine, eventoHandler *handler.EventoHandler) {
	protectedGroup := router.Group("/api/v1/eventos")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("/cria-evento", auth.RoleAuthMiddleware(auth.GetRole()), eventoHandler.CriaEvento)
		protectedGroup.GET("/lista-eventos", eventoHandler.ListaEventos)
		protectedGroup.GET("/:id", eventoHandler.GetEvento)
		protectedGroup.PUT("/atualiza-evento/:id", auth.RoleAuthMiddleware(auth.GetRole()), eventoHandler.UpdateEvento)
		protectedGroup.DELETE("/deleta-evento/:id", auth.RoleAuthMiddleware(auth.GetRole()), eventoHandler.DeleteEvento)
	}
}
