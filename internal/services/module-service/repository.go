package module_service

import (
	"context"
	"errors"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *entities.Module) error
	ReadById(ctx context.Context, id string) (*entities.Module, error)
	ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.Module, error)
	Update(ctx context.Context, module *entities.Module) (*entities.Module, error)
	Delete(ctx context.Context, id string) error
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) Create(ctx context.Context, module *entities.Module) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var position int64
	err := tx.WithContext(ctx).
		Model(&entities.Module{}).
		Where("course_id = ? AND deleted_at IS NULL", module.CourseId).
		Select("COALESCE(MAX(position),0)+1").
		Scan(&position).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	module.Position = position

	if err = tx.WithContext(ctx).Create(module).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *moduleRepository) ReadById(ctx context.Context, id string) (*entities.Module, error) {
	var module entities.Module
	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &module, nil
}

func (r *moduleRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.Module, error) {
	var modules []entities.Module
	if err := r.db.
		WithContext(ctx).
		Where("course_id = ?", courseId).
		Find(&modules).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return modules, nil
}

func (r *moduleRepository) Update(ctx context.Context, module *entities.Module) (*entities.Module, error) {
	var updatedModule entities.Module
	err := r.db.
		WithContext(ctx).
		Model(&entities.Module{}).
		Where("id = ?", module.Id).
		Updates(module).
		Scan(&updatedModule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &updatedModule, nil
}

func (r *moduleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Module{}).Error
}
