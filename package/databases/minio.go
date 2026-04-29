package databases

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client     *minio.Client
	bucketName string
}

func NewStorage(endpoint, accessKey, secretKey, bucketName string) (*Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("minio init error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket error: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("create bucket error: %w", err)
		}
	}

	return &Storage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *Storage) Upload(ctx context.Context, objectName, contentType string, size int64, reader io.Reader) error {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}

	return nil
}

func (s *Storage) GetPresignedURL(ctx context.Context, objectName string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, 15*time.Minute, nil)
	if err != nil {
		return "", fmt.Errorf("presigned url error: %w", err)
	}

	return url.String(), nil
}

func (s *Storage) Delete(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}
