package task_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entities.Task) error
	ReadById(ctx context.Context, id string) (*entities.Task, error)
	ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.Task, error)
	ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Task, error)
	Update(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Delete(ctx context.Context, id string) error
}

type taskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, task *entities.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	task.Id = uuid.NewString()

	return s.repo.Create(ctx, task)
}

func (s *taskService) ReadTaskById(ctx context.Context, id string) (*entities.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	task, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) ReadAllTasksByCourseId(ctx context.Context, courseId string) ([]entities.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tasks, err := s.repo.ReadAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *taskService) ReadAllTasksByModuleId(ctx context.Context, moduleId string) ([]entities.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tasks, err := s.repo.ReadAllByModuleId(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *taskService) UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	updatedTask, err := s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
