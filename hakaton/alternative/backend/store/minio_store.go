package store

import (
	"backend/config"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStore struct {
	client     *minio.Client
	bucketName string
	ctx        context.Context
}

func NewMinIOStore(cfg *config.MinIOConfig) (*MinIOStore, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOStore{
		client:     client,
		bucketName: cfg.BucketName,
		ctx:        ctx,
	}, nil
}

func (m *MinIOStore) UploadAvatar(userID uint64, file io.Reader, filename string, contentType string) (string, error) {
	objectName := fmt.Sprintf("avatar_user_%d", userID)

	_, err := m.client.PutObject(m.ctx, m.bucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *MinIOStore) GetAvatarURL(objectName string) string {
	return fmt.Sprintf("/api/files/%s", objectName)
}

func (m *MinIOStore) DeleteAvatar(objectName string) error {
	return m.client.RemoveObject(m.ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *MinIOStore) GetObject(objectName string) (io.Reader, error) {
	object, err := m.client.GetObject(m.ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (m *MinIOStore) GetObjectInfo(objectName string) (minio.ObjectInfo, error) {
	return m.client.StatObject(m.ctx, m.bucketName, objectName, minio.StatObjectOptions{})
}

func (m *MinIOStore) UploadObservationImage(observationID uint64, file io.Reader, filename string, contentType string) (string, error) {
	objectName := fmt.Sprintf("observation_%d_%s", observationID, filename)

	_, err := m.client.PutObject(m.ctx, m.bucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (m *MinIOStore) GetObservationImageURL(objectName string) string {
	return fmt.Sprintf("/api/files/%s", objectName)
}

func (m *MinIOStore) DeleteObservationImage(objectName string) error {
	return m.client.RemoveObject(m.ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *MinIOStore) ValidateImageType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg", 
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, allowedType := range allowedTypes {
		if strings.EqualFold(contentType, allowedType) {
			return true
		}
	}
	return false
}
