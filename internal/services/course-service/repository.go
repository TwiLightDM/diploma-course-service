package course_service

import (
	"context"
	"errors"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	ReadById(ctx context.Context, id string) (*entities.Course, error)
	ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Course, error)
	ReadAllAvailableCourses(ctx context.Context, groupsIds []string) ([]entities.Course, error)
	Update(ctx context.Context, course *entities.Course) (*entities.Course, error)
	UpdatePublishedAt(ctx context.Context, id string, time *time.Time) error
	Delete(ctx context.Context, id string) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *entities.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *courseRepository) ReadById(ctx context.Context, id string) (*entities.Course, error) {
	var course entities.Course
	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).First(&course).Error; err != nil {
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
		Where("owner_id = ?", ownerId).
		Find(&courses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) ReadAllAvailableCourses(ctx context.Context, groupsIds []string) ([]entities.Course, error) {
	var courses []entities.Course
	db := r.db.WithContext(ctx).
		Model(&entities.Course{}).
		Joins("left join group_courses on group_courses.course_id = courses.id").
		Where(
			r.db.
				Where("courses.access_type = ?", "public").
				Or(
					r.db.
						Where("courses.access_type = ?", "restricted").
						Where("group_courses.group_id IN ?", groupsIds),
				),
		).
		Distinct("courses.id")

	if err := db.Find(&courses).Error; err != nil {
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
