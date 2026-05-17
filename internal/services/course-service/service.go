package course_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	ReadById(ctx context.Context, id string) (*entities.Course, error)
	ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Course, error)
	ReadAllCourses(ctx context.Context) ([]entities.Course, error)
	ReadAllAvailableCourses(ctx context.Context, userId string) ([]entities.Course, error)
	Update(ctx context.Context, course *entities.Course) (*entities.Course, error)
	UpdatePublishedAt(ctx context.Context, id string, time *time.Time) error
	Delete(ctx context.Context, id string) error
}

type courseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(ctx context.Context, course *entities.Course) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	course.Id = uuid.NewString()

	return s.repo.Create(ctx, course)
}

func (s *courseService) ReadCourseById(ctx context.Context, id string) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	course, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) ReadAllCoursesByOwnerId(ctx context.Context, ownerId string) ([]entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	courses, err := s.repo.ReadAllByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *courseService) ReadAllCourses(ctx context.Context) ([]entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	courses, err := s.repo.ReadAllCourses(ctx)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *courseService) ReadAllAvailableCourses(ctx context.Context, userId string) ([]entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	courses, err := s.repo.ReadAllAvailableCourses(ctx, userId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *courseService) UpdateCourse(ctx context.Context, course *entities.Course) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	updatedCourse, err := s.repo.Update(ctx, course)
	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}

func (s *courseService) UpdatePublishedAt(ctx context.Context, id string) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	course, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	var t *time.Time
	if course.PublishedAt == nil {
		t = new(time.Now())
	} else {
		t = nil
	}

	err = s.repo.UpdatePublishedAt(ctx, id, t)
	if err != nil {
		return nil, err
	}

	course.PublishedAt = t

	return course, nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
