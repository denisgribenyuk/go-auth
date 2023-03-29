package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
}

var S3 S3Client

func init() {

	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: "ru-central1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	awsRegion := os.Getenv("AWS_DEFAULT_REGION")
	if awsRegion == "" {
		awsRegion = "ru-central1"
	}
	awsAccessKey := os.Getenv("S3_AUTH_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("S3_AUTH_SECRET_ACCESS_KEY")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithDefaultRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretAccessKey, "")))
	if err != nil {
		log.Fatal(err)
	}

	S3.client = s3.NewFromConfig(cfg)
}

func (c *S3Client) GetObject(bucket, key string) ([]byte, error) {
	var res []byte
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := c.client.GetObject(context.TODO(), input)
	if err != nil {
		return res, err
	}
	defer result.Body.Close()
	res, err = io.ReadAll(result.Body)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (c *S3Client) PutObject(bucket, key string, body []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}
	_, err := c.client.PutObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *S3Client) DeleteObject(bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err := c.client.DeleteObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
