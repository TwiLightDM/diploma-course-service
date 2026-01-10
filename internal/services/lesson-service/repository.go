package lesson_service

import (
	"context"
	"errors"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entities.Lesson) error
	ReadById(ctx context.Context, id string) (*entities.Lesson, error)
	ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error)
	Update(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error)
	Delete(ctx context.Context, id string) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

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
		return err
	}

	lesson.Position = position

	if err = tx.WithContext(ctx).Create(lesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *lessonRepository) ReadById(ctx context.Context, id string) (*entities.Lesson, error) {
	var lesson entities.Lesson
	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).First(&lesson).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &lesson, nil
}

func (r *lessonRepository) ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error) {
	var lessons []entities.Lesson
	if err := r.db.
		WithContext(ctx).
		Where("module_id = ?", moduleId).
		Find(&lessons).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

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

	return &updatedLesson, nil
}

func (r *lessonRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Lesson{}).Error
}
