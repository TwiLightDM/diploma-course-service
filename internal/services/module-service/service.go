package module_service

import (
	"context"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
	"time"
)

type ModuleService interface {
	CreateModule(ctx context.Context, module *entities.Module) error
	ReadModuleById(ctx context.Context, id string) (*entities.Module, error)
	ReadAllModulesByCourseId(ctx context.Context, courseId string) ([]entities.Module, error)
	UpdateModule(ctx context.Context, module *entities.Module) (*entities.Module, error)
	DeleteModule(ctx context.Context, id string) error
}

type moduleService struct {
	repo ModuleRepository
}

func NewModuleService(repo ModuleRepository) ModuleService {
	return &moduleService{repo: repo}
}

func (s *moduleService) CreateModule(ctx context.Context, module *entities.Module) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	module.Id = uuid.NewString()

	return s.repo.Create(ctx, module)
}

func (s *moduleService) ReadModuleById(ctx context.Context, id string) (*entities.Module, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	module, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (s *moduleService) ReadAllModulesByCourseId(ctx context.Context, courseId string) ([]entities.Module, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	modules, err := s.repo.ReadAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *moduleService) UpdateModule(ctx context.Context, module *entities.Module) (*entities.Module, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error

	updatedModule, err := s.repo.Update(ctx, module)
	if err != nil {
		return nil, err
	}

	return updatedModule, nil
}

func (s *moduleService) DeleteModule(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
