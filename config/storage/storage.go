package storage

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Storage *s3.Client

func New() {
	creds := credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		log.Printf("error: %v", err)
	}
	Storage = s3.NewFromConfig(cfg)
}

func Upload(key string, file *multipart.FileHeader) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(Storage)
	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	output, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   f,
	})

	return output, err
}
