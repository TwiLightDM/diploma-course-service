package lesson_progress_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
)

type LessonProgressRepository interface {
	Create(ctx context.Context, progress *entities.LessonProgress) error
	ReadByUserId(ctx context.Context, userId string) ([]entities.LessonProgress, error)
	ReadByUserIdAndLessonId(ctx context.Context, userId string, lessonId string) (*entities.LessonProgress, error)
	ReadCourseProgressByUserId(ctx context.Context, userId string, courseId string) (*entities.CourseProgress, error)
	ReadCourseStatistics(ctx context.Context, courseId string) ([]entities.UserCourseProgress, error)
	ReadModuleProgressByUserId(ctx context.Context, userId string, moduleId string) (*entities.ModuleProgress, error)
	ReadModuleStatistics(ctx context.Context, moduleId string) ([]entities.UserModuleProgress, error)
}

type lessonProgressService struct {
	repo LessonProgressRepository
}

func NewLessonProgressService(repo LessonProgressRepository) LessonProgressService {
	return &lessonProgressService{repo: repo}
}

func (s *lessonProgressService) CreateLessonProgress(ctx context.Context, progress *entities.LessonProgress) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Create(ctx, progress)
}

func (s *lessonProgressService) ReadLessonProgressByUserId(ctx context.Context, userId string) ([]entities.LessonProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progresses, err := s.repo.ReadByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return progresses, nil
}

func (s *lessonProgressService) ReadLessonProgressByUserIdAndLessonId(ctx context.Context, userId string, lessonId string) (*entities.LessonProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progress, err := s.repo.ReadByUserIdAndLessonId(ctx, userId, lessonId)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *lessonProgressService) ReadCourseProgressByUserId(ctx context.Context, userId string, courseId string) (*entities.CourseProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progress, err := s.repo.ReadCourseProgressByUserId(ctx, userId, courseId)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *lessonProgressService) ReadCourseStatistics(ctx context.Context, courseId string) ([]entities.UserCourseProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progresses, err := s.repo.ReadCourseStatistics(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return progresses, nil
}

func (s *lessonProgressService) ReadModuleProgressByUserId(ctx context.Context, userId string, moduleId string) (*entities.ModuleProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progress, err := s.repo.ReadModuleProgressByUserId(ctx, userId, moduleId)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *lessonProgressService) ReadModuleStatistics(ctx context.Context, moduleId string) ([]entities.UserModuleProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	progresses, err := s.repo.ReadModuleStatistics(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return progresses, nil
}
