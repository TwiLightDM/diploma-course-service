package completed_module_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
)

type CompletedModuleRepository interface {
	Create(ctx context.Context, completedModule *entities.CompletedModule) error
	ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedModule, error)
	ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.CompletedModule, error)
	ReadByUserIdAndModuleId(ctx context.Context, userId string, moduleId string) (*entities.CompletedModule, error)
}

type completedModuleService struct {
	repo CompletedModuleRepository
}

func NewCompletedModuleService(repo CompletedModuleRepository) CompletedModuleService {
	return &completedModuleService{repo: repo}
}

func (s *completedModuleService) CreateCompletedModule(ctx context.Context, completedModule *entities.CompletedModule) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.Create(ctx, completedModule)
}

func (s *completedModuleService) ReadAllCompletedModulesByUserId(ctx context.Context, userId string) ([]entities.CompletedModule, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedModules, err := s.repo.ReadAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return completedModules, nil
}

func (s *completedModuleService) ReadAllCompletedModulesByModuleId(ctx context.Context, moduleId string) ([]entities.CompletedModule, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedModules, err := s.repo.ReadAllByModuleId(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return completedModules, nil
}

func (s *completedModuleService) ReadCompletedModuleByUserIdAndModuleId(ctx context.Context, userId string, moduleId string) (*entities.CompletedModule, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	completedModule, err := s.repo.ReadByUserIdAndModuleId(ctx, userId, moduleId)
	if err != nil {
		return nil, err
	}

	return completedModule, nil
}
