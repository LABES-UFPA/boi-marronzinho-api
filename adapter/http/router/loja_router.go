package router

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/auth"

	"github.com/gin-gonic/gin"
)

func SetupLojaRoutes(router *gin.Engine, lojaHandler *handler.LojaHandler) {
	protectedGroup := router.Group("/api/v1/lojas")
	protectedGroup.Use(auth.JWTAuthMiddleware())
	{
		protectedGroup.GET("/todos-produtos", lojaHandler.ListaProdutos)
		protectedGroup.POST("/adiciona-produto", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.AdicionarProduto)
		protectedGroup.DELETE("/remove-produto/:id", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.RemoveProduto)
		protectedGroup.PUT("/atualiza-produto/:id", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.EditaProduto)

		protectedGroup.POST("/carrinho/adiciona-item/:id", lojaHandler.AdicionarItemCarrinho)
		protectedGroup.GET("/carrinho/:id", lojaHandler.ListarItensCarrinho)
		protectedGroup.DELETE("/carrinho/remove-item/:itemID/:usuarioID", lojaHandler.RemoverItemCarrinho)

		protectedGroup.POST("/compra/finalizar/:id", lojaHandler.FinalizarCompra)
	}
}
