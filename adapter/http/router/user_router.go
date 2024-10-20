package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	publicGroup := router.Group("/api/v1/usuarios")
	{
		publicGroup.POST("/signup", userHandler.CreateUser)
		publicGroup.POST("/login", userHandler.Login)
	}

	protectedGroup := router.Group("/api/v1/usuarios")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.DELETE("/:id", userHandler.DeleteUser)
		protectedGroup.PUT("/:id", userHandler.UpdateUser)
		protectedGroup.GET("/:id", userHandler.GetUser)
		protectedGroup.GET("/extrato/:id", userHandler.GetExtratoBoicoin)
		protectedGroup.PUT("/permissao/:id", auth.RoleAuthMiddleware(auth.GetRole()), userHandler.UpdatePermissao)
		protectedGroup.GET("todos-usuarios", auth.RoleAuthMiddleware(auth.GetRole()), userHandler.GetAllUsers)
		protectedGroup.GET("filtro", auth.RoleAuthMiddleware(auth.GetRole()), userHandler.GetUsersByName)
	}
}
