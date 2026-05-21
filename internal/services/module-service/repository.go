package module_service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type moduleRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewModuleRepository(db *gorm.DB, redis *redis.Client) ModuleRepository {
	return &moduleRepository{
		db:    db,
		redis: redis,
	}
}

const moduleCacheTTL = 10 * time.Minute

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
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {
			return errs.ErrDuplicateKey
		}
		return err
	}

	module.Position = position

	if err = tx.WithContext(ctx).Create(module).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"modules:course:"+module.CourseId,
		"course:"+module.CourseId,
		"courses:all",
	).Err()

	return nil
}

func (r *moduleRepository) ReadById(ctx context.Context, id string) (*entities.Module, error) {
	cacheKey := "module:" + id

	cachedModule, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var module entities.Module

		if json.Unmarshal([]byte(cachedModule), &module) == nil {
			return &module, nil
		}
	}

	var module entities.Module

	err = r.db.
		WithContext(ctx).
		Model(&entities.Module{}).
		Select(`
			modules.*,
			(
				SELECT COUNT(*)
				FROM lessons
				WHERE lessons.module_id = modules.id
				  AND lessons.deleted_at IS NULL
			) AS amount_of_lessons
		`).
		Where("modules.id = ?", id).
		First(&module).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	data, _ := json.Marshal(module)
	_ = r.redis.Set(ctx, cacheKey, data, moduleCacheTTL).Err()

	return &module, nil
}

func (r *moduleRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.Module, error) {
	cacheKey := "modules:course:" + courseId

	cachedModules, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var modules []entities.Module

		if json.Unmarshal([]byte(cachedModules), &modules) == nil {
			return modules, nil
		}
	}

	var modules []entities.Module

	if err = r.db.
		WithContext(ctx).
		Model(&entities.Module{}).
		Select(`
			modules.*,
			(
				SELECT COUNT(*)
				FROM lessons
				WHERE lessons.module_id = modules.id
				  AND lessons.deleted_at IS NULL
			) AS amount_of_lessons
		`).
		Where("course_id = ?", courseId).
		Find(&modules).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(modules)
	_ = r.redis.Set(ctx, cacheKey, data, moduleCacheTTL).Err()

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

	_ = r.redis.Del(ctx,
		"module:"+module.Id,
		"modules:course:"+module.CourseId,

		"course:"+module.CourseId,
		"courses:all",
	).Err()

	return &updatedModule, nil
}

func (r *moduleRepository) Delete(ctx context.Context, id string) error {
	module, err := r.ReadById(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&entities.Module{}).Error

	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"module:"+id,
		"modules:course:"+module.CourseId,

		"course:"+module.CourseId,
		"courses:all",
	).Err()

	return nil
}
