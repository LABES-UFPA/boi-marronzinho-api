package container

import (
	"boi-marronzinho-api/postgres"

	"gorm.io/gorm"
)

type Container struct {
	DB          *gorm.DB
}

func NewContainer() *Container {
	db := postgres.InitDB()

	return &Container{
		DB:          db,
	}
}
