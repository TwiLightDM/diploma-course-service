package lesson_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/package/databases"
)

type lessonStorage struct {
	storage *databases.Storage
}

func NewLessonStorage(storage *databases.Storage) LessonStorage {
	return &lessonStorage{storage: storage}
}

func (s *lessonStorage) ReadAllByLessonId(ctx context.Context, lessons []entities.Lesson) error {
	for i := range lessons {
		for j := range lessons[i].Files {

			url, err := s.storage.GetPresignedURL(ctx, lessons[i].Files[j].ObjectName)
			if err != nil {
				return err
			}

			lessons[i].Files[j].Url = url
		}
	}

	return nil
}

func (s *lessonStorage) ReadByObjectName(ctx context.Context, objectName string) (string, error) {
	url, err := s.storage.GetPresignedURL(ctx, objectName)
	if err != nil {
		return "", err
	}

	return url, nil
}
