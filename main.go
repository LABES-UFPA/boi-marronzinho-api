package main

import (
	"boi-marronzinho-api/adapter/http"
	"boi-marronzinho-api/core/usecase"
	"boi-marronzinho-api/minio"
	"boi-marronzinho-api/module/boicoin"
	"boi-marronzinho-api/module/evento"
	"boi-marronzinho-api/module/loja"
	"boi-marronzinho-api/module/oficina"
	"boi-marronzinho-api/module/troca"
	"boi-marronzinho-api/module/user"
	"boi-marronzinho-api/postgres"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(postgres.InitDB),      
		fx.Provide(minio.InitMinio),      
		http.RouterModule(),              
		user.UserModule,                  
		troca.TrocaModule,                
		boicoin.BoicoinModule,            
		oficina.OficnaModule,             
		loja.LojaModule,                  
		evento.EventoModule,              
		fx.Invoke(http.RegisterRoutes),   
		fx.Invoke(iniciarCronJobs),
	)

	app.Run()
}

func iniciarCronJobs(tuc *usecase.TrocaUseCase) {
	tuc.IniciarCronJobExpiracaoTroca()
}
