package completed_course_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"

	"gorm.io/gorm"
)

type completedCourseRepository struct {
	db *gorm.DB
}

func NewCompletedCourseRepository(db *gorm.DB) CompletedCourseRepository {
	return &completedCourseRepository{db: db}
}

func (r *completedCourseRepository) Create(ctx context.Context, completedCourse *entities.CompletedCourse) error {
	return r.db.WithContext(ctx).Create(completedCourse).Error
}

func (r *completedCourseRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.CompletedCourse, error) {
	var completedCourses []entities.CompletedCourse
	if err := r.db.
		WithContext(ctx).
		Where("course_id = ?", courseId).
		Find(&completedCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedCourses, nil
}

func (r *completedCourseRepository) ReadAllByUserId(ctx context.Context, userId string) ([]entities.CompletedCourse, error) {
	var completedCourses []entities.CompletedCourse
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&completedCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return completedCourses, nil
}

func (r *completedCourseRepository) ReadByUserIdAndCourseId(ctx context.Context, userId, courseId string) (*entities.CompletedCourse, error) {
	var completedCourse entities.CompletedCourse
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? and course_id = ?", userId, courseId).
		First(&completedCourse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &completedCourse, nil
}
