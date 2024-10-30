package minio

import (
    "context"
    "fmt"
    "io"
    "net/http"

    "github.com/minio/minio-go/v7"
)

func UploadFile(file io.Reader, fileName string, bucketName string) (string, error) {
    if MinioClient == nil {
        return "", fmt.Errorf("MinioClient não está inicializado")
    }

    contentType, err := detectContentType(file)
    if err != nil {
        return "", fmt.Errorf("erro ao detectar tipo de conteúdo: %v", err)
    }

    exists, err := MinioClient.BucketExists(context.Background(), bucketName)
    if err != nil {
        return "", fmt.Errorf("erro ao verificar o bucket: %v", err)
    }
    if !exists {
        err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return "", fmt.Errorf("erro ao criar bucket: %v", err)
        }
    }

    _, err = MinioClient.PutObject(context.Background(), bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
    if err != nil {
        return "", fmt.Errorf("erro ao fazer upload do arquivo: %v", err)
    }

    return fmt.Sprintf("/%s/%s", bucketName, fileName), nil
}

func DeleteFile(fileName string, bucketName string) error {
    if MinioClient == nil {
        return fmt.Errorf("MinioClient não está inicializado")
    }

    err := MinioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
    if err != nil {
        return fmt.Errorf("erro ao deletar arquivo: %v", err)
    }

    return nil
}

func UpdateFile(file io.Reader, fileName string, bucketName string) (string, error) {
    if MinioClient == nil {
        return "", fmt.Errorf("MinioClient não está inicializado")
    }

    _ = DeleteFile(fileName, bucketName)

    contentType, err := detectContentType(file)
    if err != nil {
        return "", fmt.Errorf("erro ao detectar tipo de conteúdo: %v", err)
    }

    _, err = MinioClient.PutObject(context.Background(), bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
    if err != nil {
        return "", fmt.Errorf("erro ao atualizar arquivo: %v", err)
    }

    return fmt.Sprintf("/%s/%s", bucketName, fileName), nil
}


func detectContentType(file io.Reader) (string, error) {
    // Lê os primeiros 512 bytes para detectar o tipo de conteúdo
    buffer := make([]byte, 512)
    _, err := file.Read(buffer)
    if err != nil {
        return "", err
    }

    contentType := http.DetectContentType(buffer)

    if seeker, ok := file.(io.Seeker); ok {
        seeker.Seek(0, io.SeekStart)
    }

    return contentType, nil
}
