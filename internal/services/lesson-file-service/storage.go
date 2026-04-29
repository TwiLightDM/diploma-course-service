package lesson_file_service

import (
	"bytes"
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/package/databases"
	"golang.org/x/sync/errgroup"
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
	urls := make([]string, len(lessonFiles))
	g, ctx := errgroup.WithContext(ctx)

	for i := range lessonFiles {
		g.Go(func(i int) func() error {
			return func() error {
				url, err := s.storage.GetPresignedURL(ctx, lessonFiles[i].ObjectName)
				if err != nil {
					return err
				}

				lessonFiles[i].Url = url
				return nil
			}
		}(i))
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for i := range lessonFiles {
		lessonFiles[i].Url = urls[i]
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
