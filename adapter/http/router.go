package http

import (
	"boi-marronzinho-api/adapter/http/handler"
	"boi-marronzinho-api/adapter/http/router"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func StartServer(lc fx.Lifecycle, r *gin.Engine) {
	port := ":8080"
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logrus.Infof("API rodando na porta %s", port)
				if err := r.Run(port); err != nil {
					logrus.Errorf("Erro ao iniciar o servidor: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Info("API foi finalizada")
			return nil
		},
	})
}

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	doacaoHandler *handler.DoacaoHandler,
	oficinaHandler *handler.OficinaHandler,
	carteiraHandler *handler.CarteiraHandler,
) {
	router.SetupUserRoutes(r, userHandler)
	router.SetupDoacaoRoutes(r, doacaoHandler)
	router.SetupOficinaRoutes(r, oficinaHandler)
	router.SetupCarteiraRoutes(r, carteiraHandler)
}

func SetupRouter() *gin.Engine {
	return gin.Default()
}

func RouterModule() fx.Option {
	return fx.Options(
		fx.Provide(SetupRouter),
		fx.Invoke(StartServer),
	)
}
