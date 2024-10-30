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
		// Rotas de Produto
		protectedGroup.GET("/todos-produtos", lojaHandler.ListaProdutos)
		protectedGroup.POST("/adiciona-produto", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.AdicionarProduto)
		protectedGroup.DELETE("/remove-produto/:id", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.RemoveProduto)
		protectedGroup.PUT("/atualiza-produto/:id", auth.RoleAuthMiddleware(auth.GetRole()), lojaHandler.EditaProduto)

		// Rotas do Carrinho
		protectedGroup.POST("/carrinho/adiciona-item/:usuarioID", lojaHandler.AdicionarItemCarrinho)
		protectedGroup.GET("/carrinho/:usuarioID", lojaHandler.ListarItensCarrinho)
		protectedGroup.DELETE("/carrinho/remove-item/:itemID/:usuarioID", lojaHandler.RemoverItemCarrinho)

		// Novo Endpoint: Atualizar ou Remover Unidade de Item do Carrinho
		protectedGroup.PUT("/carrinho/atualiza-quantidade-item/:usuarioID", lojaHandler.AtualizarItemCarrinho)

		// Rota para Finalizar Compra
		protectedGroup.POST("/compra/finalizar/:usuarioID", lojaHandler.FinalizarCompra)
	}
}
