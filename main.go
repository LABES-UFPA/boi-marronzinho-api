package main

import (
	"boi-marronzinho-api/adapter/http"
	"boi-marronzinho-api/module/user"
	"boi-marronzinho-api/postgres"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(postgres.InitDB),
		http.RouterModule(),
		user.UserModule,
		fx.Invoke(http.RegisterRoutes),
	)

	app.Run()
}
