package lesson_file_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"gorm.io/gorm"
)

type lessonFileRepository struct {
	db *gorm.DB
}

func NewLessonFileRepository(db *gorm.DB) LessonFileRepository {
	return &lessonFileRepository{db: db}
}

func (r *lessonFileRepository) Create(ctx context.Context, lessonFile *entities.LessonFile) error {
	err := r.db.WithContext(ctx).Create(lessonFile).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *lessonFileRepository) ReadAllByLessonId(ctx context.Context, lessonId string) ([]entities.LessonFile, error) {
	var lessonFiles []entities.LessonFile
	if err := r.db.
		WithContext(ctx).
		Where("lesson_id = ?", lessonId).
		Find(&lessonFiles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return lessonFiles, nil
}

func (r *lessonFileRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.LessonFile{}).Error
}
