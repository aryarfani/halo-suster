package storage

import (
	"os"

	"github.com/gofiber/storage/s3/v2"
)

var Storage *s3.Storage

func New() {
	Storage = s3.New(s3.Config{
		Bucket: os.Getenv("AWS_S3_BUCKET_NAME"),
		Region: os.Getenv("AWS_REGION"),
		Credentials: s3.Credentials{
			AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})
}
