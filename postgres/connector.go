package postgres

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error

	err = godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("Connecting to PostgreSQL...")

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" TimeZone=" + os.Getenv("DB_TIMEZONE") +
		" connect_timeout=" + os.Getenv("DB_CONNECT_TIMEOUT")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":     os.Getenv("DB_HOST"),
			"user":     os.Getenv("DB_USER"),
			"database": os.Getenv("DB_NAME"),
		}).WithError(err).Fatal("Failed to connect to database")
	}

	logrus.Info("Database connected successfully")
	return DB
}
