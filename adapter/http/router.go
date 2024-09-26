package http

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/http/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UsuarioHandler,
) {
	router.SetupUsuarioRoutes(r, userHandler)
}

func SetupRouter() *gin.Engine {
	return gin.Default()
}

func RouterModule() fx.Option {
	return fx.Provide(SetupRouter)
}
