package lesson_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entities.Lesson) error
	ReadById(ctx context.Context, id string) (*entities.Lesson, error)
	ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error)
	Update(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error)
	Delete(ctx context.Context, id string) error
}

type LessonStorage interface {
	ReadAllByLessonId(ctx context.Context, lessons []entities.Lesson) error
	ReadByObjectName(ctx context.Context, objectName string) (string, error)
}

type lessonService struct {
	repo  LessonRepository
	store LessonStorage
}

func NewLessonService(repo LessonRepository, store LessonStorage) LessonService {
	return &lessonService{repo: repo, store: store}
}

func (s *lessonService) CreateLesson(ctx context.Context, lesson *entities.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lesson.Id = uuid.NewString()

	return s.repo.Create(ctx, lesson)
}

func (s *lessonService) ReadLessonById(ctx context.Context, id string) (*entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lesson, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range lesson.Files {
		lesson.Files[i].Url, err = s.store.ReadByObjectName(ctx, lesson.Files[i].ObjectName)
		if err != nil {
			return nil, err
		}
	}

	return lesson, nil
}

func (s *lessonService) ReadAllLessonsByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lessons, err := s.repo.ReadAllByModuleId(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	err = s.store.ReadAllByLessonId(ctx, lessons)

	return lessons, nil
}

func (s *lessonService) UpdateLesson(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	updatedLesson, err := s.repo.Update(ctx, lesson)
	if err != nil {
		return nil, err
	}

	for i := range lesson.Files {
		lesson.Files[i].Url, err = s.store.ReadByObjectName(ctx, lesson.Files[i].ObjectName)
		if err != nil {
			return nil, err
		}
	}

	return updatedLesson, nil
}

func (s *lessonService) DeleteLesson(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
