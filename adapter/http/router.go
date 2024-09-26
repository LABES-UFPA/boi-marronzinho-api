package http

import (
	"boi-marronzinho-api/adapter/http/router"
	"boi-marronzinho-api/container"

	"github.com/gin-gonic/gin"
)


func SetupRouter(cont *container.Container) *gin.Engine {
	r := gin.Default()

	router.SetupUsuarioRoutes(r, cont.UserHandler)

	return r
}