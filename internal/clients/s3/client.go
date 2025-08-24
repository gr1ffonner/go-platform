package s3

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type clientS3 struct {
	client             *s3.Client
	bucketName         string
	baseEndpoint       string
	basePublicEndpoint string
}

// NewClientS3 создает новый экземпляр клиента
func NewClientS3(keyID, keySecret, bucketName, baseEndpoint, basePublicEndpoint, region string) (*clientS3, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(keyID, keySecret, ""),
		),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(baseEndpoint)
		o.UsePathStyle = true
	})

	return &clientS3{
		client:             client,
		bucketName:         bucketName,
		baseEndpoint:       baseEndpoint,
		basePublicEndpoint: basePublicEndpoint,
	}, nil
}

// PutObject — загрузить объект в бакет
func (c *clientS3) PutObject(ctx context.Context, key string, data []byte) (err error) {
	input := &s3.PutObjectInput{
		Bucket: &c.bucketName,
		Key:    &key,
		Body:   bytes.NewReader(data),
	}
	_, err = c.client.PutObject(ctx, input)
	if err != nil {
		return
	}

	return
}

// GenerateURL generates a public URL for the given key
func (c *clientS3) GenerateURL(key string) string {
	return strings.Join([]string{c.basePublicEndpoint, c.bucketName, key}, "/")
}
