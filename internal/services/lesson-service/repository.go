package lesson_service

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

type lessonRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewLessonRepository(db *gorm.DB, redis *redis.Client) LessonRepository {
	return &lessonRepository{
		db:    db,
		redis: redis,
	}
}

const lessonCacheTTL = 10 * time.Minute

func (r *lessonRepository) Create(ctx context.Context, lesson *entities.Lesson) error {
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
		Model(&entities.Lesson{}).
		Where("module_id = ? AND deleted_at IS NULL", lesson.ModuleId).
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

	lesson.Position = position

	if err = tx.WithContext(ctx).Create(lesson).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	_ = r.redis.Del(ctx, "lessons:module:"+lesson.ModuleId).Err()

	return nil
}

func (r *lessonRepository) ReadById(ctx context.Context, id string) (*entities.Lesson, error) {
	cacheKey := "lesson:" + id

	cachedLesson, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var lesson entities.Lesson

		if json.Unmarshal([]byte(cachedLesson), &lesson) == nil {
			return &lesson, nil
		}
	}

	var lesson entities.Lesson

	if err = r.db.
		WithContext(ctx).
		Preload("Files").
		Where("id = ?", id).
		First(&lesson).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	data, _ := json.Marshal(lesson)
	_ = r.redis.Set(ctx, cacheKey, data, lessonCacheTTL).Err()

	return &lesson, nil
}

func (r *lessonRepository) ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error) {
	cacheKey := "lessons:module:" + moduleId

	cachedLessons, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var lessons []entities.Lesson

		if json.Unmarshal([]byte(cachedLessons), &lessons) == nil {
			return lessons, nil
		}
	}

	var lessons []entities.Lesson

	if err = r.db.
		WithContext(ctx).
		Preload("Files").
		Where("module_id = ?", moduleId).
		Find(&lessons).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(lessons)
	_ = r.redis.Set(ctx, cacheKey, data, lessonCacheTTL).Err()

	return lessons, nil
}

func (r *lessonRepository) Update(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error) {
	var updatedLesson entities.Lesson

	err := r.db.
		WithContext(ctx).
		Model(&entities.Lesson{}).
		Where("id = ?", lesson.Id).
		Updates(lesson).
		Scan(&updatedLesson).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	_ = r.redis.Del(ctx,
		"lesson:"+lesson.Id,
		"lessons:module:"+lesson.ModuleId,
	).Err()

	return &updatedLesson, nil
}

func (r *lessonRepository) Delete(ctx context.Context, id string) error {
	lesson, err := r.ReadById(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&entities.Lesson{}).Error

	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"lesson:"+id,
		"lessons:module:"+lesson.ModuleId,
	).Err()

	return nil
}
