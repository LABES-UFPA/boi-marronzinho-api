package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

var MinioClient *minio.Client

func InitMinio() (*minio.Client, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("Tentando inicializar o cliente MinIO...")

	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "password", ""),
		Secure: false,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Erro ao conectar ao MinIO")
		return nil, err
	}

	logrus.Info("Cliente MinIO inicializado com sucesso")

	MinioClient = client
	return client, nil
}
