package completed_theory_course_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
)

type CompletedTheoryCourseRepository interface {
	Create(ctx context.Context, completedTheoryCourse *entities.CompletedTheoryCourse) error
	ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedTheoryCourse, error)
	ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.CompletedTheoryCourse, error)
	ReadByUserIdAndCourseId(ctx context.Context, userId string, courseId string) (*entities.CompletedTheoryCourse, error)
}

type completedTheoryCourseService struct {
	repo CompletedTheoryCourseRepository
}

func NewCompletedTheoryCourseService(repo CompletedTheoryCourseRepository) CompletedTheoryCourseService {
	return &completedTheoryCourseService{repo: repo}
}

func (s *completedTheoryCourseService) CreateCompletedTheoryCourse(ctx context.Context, completedTheoryCourse *entities.CompletedTheoryCourse) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Create(ctx, completedTheoryCourse)
}

func (s *completedTheoryCourseService) ReadAllCompletedTheoryCoursesByUserId(ctx context.Context, userId string) ([]entities.CompletedTheoryCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedTheoryCourses, err := s.repo.ReadAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return completedTheoryCourses, nil
}

func (s *completedTheoryCourseService) ReadAllCompletedTheoryCoursesByCourseId(ctx context.Context, courseId string) ([]entities.CompletedTheoryCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedTheoryCourses, err := s.repo.ReadAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return completedTheoryCourses, nil
}

func (s *completedTheoryCourseService) ReadCompletedTheoryCourseByUserIdAndCourseId(ctx context.Context, userId string, courseId string) (*entities.CompletedTheoryCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedTheoryCourse, err := s.repo.ReadByUserIdAndCourseId(ctx, userId, courseId)
	if err != nil {
		return nil, err
	}

	return completedTheoryCourse, nil
}
