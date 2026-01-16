package lesson_service

import (
	"context"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
	"time"
)

type LessonService interface {
	CreateLesson(ctx context.Context, lesson *entities.Lesson) error
	ReadLessonById(ctx context.Context, id string) (*entities.Lesson, error)
	ReadAllLessonsByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
}

type lessonService struct {
	repo LessonRepository
}

func NewLessonService(repo LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) CreateLesson(ctx context.Context, lesson *entities.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	lesson.Id = uuid.NewString()

	return s.repo.Create(ctx, lesson)
}

func (s *lessonService) ReadLessonById(ctx context.Context, id string) (*entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	lesson, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return lesson, nil
}

func (s *lessonService) ReadAllLessonsByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	lessons, err := s.repo.ReadAllByModuleId(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (s *lessonService) UpdateLesson(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error

	updatedLesson, err := s.repo.Update(ctx, lesson)
	if err != nil {
		return nil, err
	}

	return updatedLesson, nil
}

func (s *lessonService) DeleteLesson(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
