package completed_theory_course_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"

	"gorm.io/gorm"
)

type completedTheoryCourseRepository struct {
	db *gorm.DB
}

func NewCompletedTheoryCourseRepository(db *gorm.DB) CompletedTheoryCourseRepository {
	return &completedTheoryCourseRepository{db: db}
}

func (r *completedTheoryCourseRepository) Create(ctx context.Context, completedTheoryCourse *entities.CompletedTheoryCourse) error {
	return r.db.WithContext(ctx).Create(completedTheoryCourse).Error
}

func (r *completedTheoryCourseRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.CompletedTheoryCourse, error) {
	var completedTheoryCourses []entities.CompletedTheoryCourse
	if err := r.db.
		WithContext(ctx).
		Where("course_id = ?", courseId).
		Find(&completedTheoryCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedTheoryCourses, nil
}

func (r *completedTheoryCourseRepository) ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedTheoryCourse, error) {
	var completedTheoryCourses []entities.CompletedTheoryCourse
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&completedTheoryCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedTheoryCourses, nil
}

func (r *completedTheoryCourseRepository) ReadByUserIdAndCourseId(ctx context.Context, userId, courseId string) (*entities.CompletedTheoryCourse, error) {
	var completedTheoryCourse entities.CompletedTheoryCourse
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? and course_id = ?", userId, courseId).
		First(&completedTheoryCourse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &completedTheoryCourse, nil
}
