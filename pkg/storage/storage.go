/*
ZAU Thoth API
Copyright (C) 2021 Daniel A. Hawton (daniel@hawton.org)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package storage

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

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

func GetFileStream(storageUrl string) (*minio.Object, error) {
	u, err := url.Parse(storageUrl)
	if err != nil {
		return nil, err
	}
	var minioClient *minio.Client
	var bucket string

	// Allow for multiple buckets, so check the host here
	if u.Host == "do.chicagoartcc.org" {
		minioClient, err = prepConnection()
		bucket = "vzau"
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported storage provider")
	}

	key := strings.TrimLeft(u.Path, "/")
	file, err := minioClient.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return file, nil
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
