package lesson_file_service

import (
	"context"
	"fmt"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type LessonFileRepository interface {
	Create(ctx context.Context, lessonFile *entities.LessonFile) error
	ReadAllByLessonId(ctx context.Context, lessonId string) ([]entities.LessonFile, error)
	Delete(ctx context.Context, id string) error
}

type LessonFileStorage interface {
	Upload(ctx context.Context, lessonFile *entities.LessonFile, file []byte) error
	ReadAllByLessonId(ctx context.Context, lessonFiles []entities.LessonFile) error
	ReadByObjectName(ctx context.Context, objectName string) (string, error)
	Delete(ctx context.Context, objectName string) error
}

type lessonFileService struct {
	repo  LessonFileRepository
	store LessonFileStorage
}

func NewLessonFileService(repo LessonFileRepository, store LessonFileStorage) LessonFileService {
	return &lessonFileService{repo: repo, store: store}
}

func (s *lessonFileService) UploadFile(ctx context.Context, lessonFile *entities.LessonFile, file []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lessonFile.Id = uuid.NewString()
	lessonFile.ObjectName = fmt.Sprintf("lessons/%s/%s/%s", lessonFile.LessonId, lessonFile.Id, lessonFile.FileName)
	lessonFile.Size = int64(len(file))

	err := s.repo.Create(ctx, lessonFile)
	if err != nil {
		return err
	}

	err = s.store.Upload(ctx, lessonFile, file)
	if err != nil {
		return err
	}

	lessonFile.Url, err = s.store.ReadByObjectName(ctx, lessonFile.ObjectName)

	return nil
}

func (s *lessonFileService) GetLessonFiles(ctx context.Context, lessonId string) ([]entities.LessonFile, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lessonFiles, err := s.repo.ReadAllByLessonId(ctx, lessonId)
	if err != nil {
		return nil, err
	}

	err = s.store.ReadAllByLessonId(ctx, lessonFiles)
	if err != nil {
		return nil, err
	}

	return lessonFiles, nil
}

func (s *lessonFileService) DeleteFile(ctx context.Context, id, objectName string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.Delete(ctx, objectName)
	if err != nil {
		return err
	}

	return nil
}
