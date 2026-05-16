package completed_course_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
)

type CompletedCourseRepository interface {
	Create(ctx context.Context, completedCourse *entities.CompletedCourse) error
	ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedCourse, error)
	ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.CompletedCourse, error)
	ReadByUserIdAndCourseId(ctx context.Context, userId string, courseId string) (*entities.CompletedCourse, error)
}

type completedCourseService struct {
	repo CompletedCourseRepository
}

func NewCompletedCourseService(repo CompletedCourseRepository) CompletedCourseService {
	return &completedCourseService{repo: repo}
}

func (s *completedCourseService) CreateCompletedCourse(ctx context.Context, completedCourse *entities.CompletedCourse) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Create(ctx, completedCourse)
}

func (s *completedCourseService) ReadAllCompletedCoursesByUserId(ctx context.Context, userId string) ([]entities.CompletedCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedCourses, err := s.repo.ReadAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return completedCourses, nil
}

func (s *completedCourseService) ReadAllCompletedCoursesByCourseId(ctx context.Context, courseId string) ([]entities.CompletedCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedCourses, err := s.repo.ReadAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return completedCourses, nil
}

func (s *completedCourseService) ReadCompletedCourseByUserIdAndCourseId(ctx context.Context, userId string, courseId string) (*entities.CompletedCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedCourse, err := s.repo.ReadByUserIdAndCourseId(ctx, userId, courseId)
	if err != nil {
		return nil, err
	}

	return completedCourse, nil
}
