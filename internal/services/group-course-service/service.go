package group_course_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type GroupCourseService interface {
	CreateGroupCourse(ctx context.Context, groupCourse *entities.GroupCourse) (*entities.GroupCourse, error)
	ReadAllGroupCoursesByCourseId(ctx context.Context, courseId string) ([]entities.GroupCourse, error)
	ReadAllGroupCoursesByGroupId(ctx context.Context, groupId string) ([]entities.GroupCourse, error)
	DeleteGroupCourse(ctx context.Context, id string) error
}

type groupCourseService struct {
	repo GroupCourseRepository
}

func NewGroupCourseService(repo GroupCourseRepository) GroupCourseService {
	return &groupCourseService{repo: repo}
}

func (s *groupCourseService) CreateGroupCourse(ctx context.Context, groupCourse *entities.GroupCourse) (*entities.GroupCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupCourse.Id = uuid.NewString()

	err := s.repo.Create(ctx, groupCourse)
	if err != nil {
		return nil, err
	}

	return groupCourse, nil
}

func (s *groupCourseService) ReadAllGroupCoursesByCourseId(ctx context.Context, courseId string) ([]entities.GroupCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupCourse, err := s.repo.ReadAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return groupCourse, nil
}

func (s *groupCourseService) ReadAllGroupCoursesByGroupId(ctx context.Context, groupId string) ([]entities.GroupCourse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupCourses, err := s.repo.ReadAllByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	return groupCourses, nil
}

func (s *groupCourseService) DeleteGroupCourse(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
