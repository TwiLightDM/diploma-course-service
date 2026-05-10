package course_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *entities.Course) error {
	err := r.db.WithContext(ctx).Create(course).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {
			return errs.ErrDuplicateKey
		}
	}

	return nil
}

func (r *courseRepository) ReadById(ctx context.Context, id string) (*entities.Course, error) {
	var course entities.Course

	err := r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Select(`
			courses.*,
			(
				SELECT COUNT(*)
				FROM modules
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
			) AS amount_of_modules,
			(
				SELECT COUNT(*)
				FROM lessons
				JOIN modules ON modules.id = lessons.module_id
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
				  AND lessons.deleted_at IS NULL
			) AS amount_of_lessons
		`).
		Where("courses.id = ?", id).
		First(&course).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Course, error) {
	var courses []entities.Course
	if err := r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Select(`
			courses.*,
			(
				SELECT COUNT(*)
				FROM modules
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
			) AS amount_of_modules,
			(
				SELECT COUNT(*)
				FROM lessons
				JOIN modules ON modules.id = lessons.module_id
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
				  AND lessons.deleted_at IS NULL
			) AS amount_of_lessons
		`).
		Where("owner_id = ?", ownerId).
		Find(&courses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) ReadAllCourses(ctx context.Context) ([]entities.Course, error) {
	var courses []entities.Course
	if err := r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Select(`
			courses.*,
			(
				SELECT COUNT(*)
				FROM modules
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
			) AS amount_of_modules,
			(
				SELECT COUNT(*)
				FROM lessons
				JOIN modules ON modules.id = lessons.module_id
				WHERE modules.course_id = courses.id
				  AND modules.deleted_at IS NULL
				  AND lessons.deleted_at IS NULL
			) AS amount_of_lessons
		`).
		Find(&courses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) Update(ctx context.Context, course *entities.Course) (*entities.Course, error) {
	var updatedCourse entities.Course
	err := r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Where("id = ?", course.Id).
		Updates(course).
		Scan(&updatedCourse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &updatedCourse, nil
}

func (r *courseRepository) UpdatePublishedAt(ctx context.Context, id string, time *time.Time) error {
	return r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"published_at": time,
		}).Error
}

func (r *courseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Course{}).Error
}
