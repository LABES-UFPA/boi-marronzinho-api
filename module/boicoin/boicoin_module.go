package boicoin

import (
	"boi-marronzinho-api/adapter/repository"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var BoincoinModule = fx.Options(
    fx.Provide(NewBoicoinRepository),
)

func NewBoicoinRepository(db *gorm.DB) repository.BoicoinRepository {
    return repository.NewBoicoinRepository(db)
}