package router

import (
	"boi-marronzinho-api/adapter/http/handler"

	"github.com/gin-gonic/gin"
)


func SetupUsuarioRoutes(router *gin.Engine, userHandler *handler.UsuarioHandler) {
	userGroup := router.Group("/api/v1/usuarios")
	{
		userGroup.POST("", userHandler.CreateUser)
		// userGroup.POST("/login", userHandler.Login)
		// userGroup.GET("/:id", userHandler.GetUser)
	}
}
