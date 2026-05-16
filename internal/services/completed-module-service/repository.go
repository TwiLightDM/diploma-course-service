package completed_module_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"gorm.io/gorm"
)

type completedModuleRepository struct {
	db *gorm.DB
}

func NewCompletedModuleRepository(db *gorm.DB) CompletedModuleRepository {
	return &completedModuleRepository{db: db}
}

func (r *completedModuleRepository) Create(ctx context.Context, completedModule *entities.CompletedModule) error {
	return r.db.WithContext(ctx).Create(completedModule).Error
}

func (r *completedModuleRepository) ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.CompletedModule, error) {
	var completedModules []entities.CompletedModule
	if err := r.db.
		WithContext(ctx).
		Where("module_id = ?", moduleId).
		Find(&completedModules).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedModules, nil
}

func (r *completedModuleRepository) ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedModule, error) {
	var completedModules []entities.CompletedModule
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&completedModules).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedModules, nil
}

func (r *completedModuleRepository) ReadByUserIdAndModuleId(ctx context.Context, userId, moduleId string) (*entities.CompletedModule, error) {
	var completedModule entities.CompletedModule
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? and module_id = ?", userId, moduleId).
		First(&completedModule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &completedModule, nil
}
