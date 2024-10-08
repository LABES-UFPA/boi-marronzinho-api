package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"
	"boi-marronzinho-api/global/enums"

	"github.com/gin-gonic/gin"
)

func SetupOficinaRoutes(router *gin.Engine, oficinaHandler *handler.OficinaHandler) {
	protectedGroup := router.Group("/api/v1/oficinas")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.POST("/inscricao", oficinaHandler.InscricaoOficina)
		protectedGroup.GET("/lista-oficinas", oficinaHandler.ListaOficinas)
		protectedGroup.POST("/cria-oficinas", auth.RoleAuthMiddleware(getRole()), oficinaHandler.CriaOficina)
	}
}

func getRole() string {
	role, err := enums.GetUserRole(1)
	if err != nil {
		return "Role desconhecida"
	}
	return role
}
