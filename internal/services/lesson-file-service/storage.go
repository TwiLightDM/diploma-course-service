package lesson_file_service

import (
	"bytes"
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/package/databases"
)

type lessonFileStorage struct {
	storage *databases.Storage
}

func NewLessonFileStorage(storage *databases.Storage) LessonFileStorage {
	return &lessonFileStorage{storage: storage}
}

func (s *lessonFileStorage) Upload(ctx context.Context, lessonFile *entities.LessonFile, file []byte) error {
	err := s.storage.Upload(
		ctx,
		lessonFile.ObjectName,
		lessonFile.ContentType,
		lessonFile.Size,
		bytes.NewReader(file),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *lessonFileStorage) ReadAllByLessonId(ctx context.Context, lessonFiles []entities.LessonFile) error {
	var err error
	for i := range lessonFiles {
		lessonFiles[i].Url, err = s.storage.GetPresignedURL(ctx, lessonFiles[i].ObjectName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *lessonFileStorage) ReadByObjectName(ctx context.Context, objectName string) (string, error) {
	url, err := s.storage.GetPresignedURL(ctx, objectName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *lessonFileStorage) Delete(ctx context.Context, objectName string) error {
	return s.storage.Delete(ctx, objectName)
}
