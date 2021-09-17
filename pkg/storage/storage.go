package storage

import (
	"context"
	"net/http"
	"os"

	"github.com/dhawton/log4g"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vzau/common/utils"
)

var log = log4g.Category("s3")

func prepConnection() (*minio.Client, error) {
	endpoint := utils.Getenv("AWS_ENDPOINT", "")
	accessKey := utils.Getenv("AWS_ACCESS_KEY", "")
	secretKey := utils.Getenv("AWS_SECRET_KEY", "")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Error("Error creating minio client: %s", err.Error())
		return nil, err
	}

	log.Debug("Created minio client")
	return minioClient, nil
}

func GetContentType(filePath string) string {
	f, _ := os.Open(filePath)
	defer f.Close()

	buffer := make([]byte, 512)
	f.Read(buffer)
	return http.DetectContentType(buffer)
}

func UploadFile(bucket string, key string, filePath string, contentType string) error {
	minioClient, err := prepConnection()
	if err != nil {
		return err
	}

	_, err = minioClient.FPutObject(context.Background(), bucket, key, filePath, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		log.Error("Error uploading file: %s", err.Error())
		return err
	}

	log.Debug("Uploaded file to storage")
	return nil
}

func DeleteFile(bucket string, key string) error {
	minioClient, err := prepConnection()
	if err != nil {
		return err
	}

	err = minioClient.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		log.Error("Error deleting file: %s", err.Error())
		return err
	}

	log.Debug("Deleted file from storage")
	return nil
}
