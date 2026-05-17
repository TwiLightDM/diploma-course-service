package course_service

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

type courseRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCourseRepository(db *gorm.DB, redis *redis.Client) CourseRepository {
	return &courseRepository{
		db:    db,
		redis: redis,
	}
}

const courseCacheTTL = 10 * time.Minute

func (r *courseRepository) Create(ctx context.Context, course *entities.Course) error {
	err := r.db.WithContext(ctx).Create(course).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {
			return errs.ErrDuplicateKey
		}

		return err
	}

	_ = r.redis.Del(ctx, "courses:all").Err()

	return nil
}

func (r *courseRepository) ReadById(ctx context.Context, id string) (*entities.Course, error) {
	cacheKey := "course:" + id

	cachedCourse, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var course entities.Course

		if json.Unmarshal([]byte(cachedCourse), &course) == nil {
			return &course, nil
		}
	}

	var course entities.Course

	err = r.db.
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

	data, _ := json.Marshal(course)
	_ = r.redis.Set(ctx, cacheKey, data, courseCacheTTL).Err()

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
		Order("published_at IS NOT NULL").
		Order("published_at DESC").
		Find(&courses).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) ReadAllCourses(ctx context.Context) ([]entities.Course, error) {
	const cacheKey = "courses:all"

	cachedCourses, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var courses []entities.Course

		if json.Unmarshal([]byte(cachedCourses), &courses) == nil {
			return courses, nil
		}
	}

	var courses []entities.Course

	if err = r.db.
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

	data, _ := json.Marshal(courses)
	_ = r.redis.Set(ctx, cacheKey, data, courseCacheTTL).Err()

	return courses, nil
}

func (r *courseRepository) ReadAllAvailableCourses(ctx context.Context, userId string) ([]entities.Course, error) {
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
		Where(`
			courses.published_at IS NOT NULL
			AND (
				courses.access_type = 'public'
				OR (
					courses.access_type = 'group_only'
					AND EXISTS (
						SELECT 1
						FROM group_courses gc
						JOIN group_members gm
							ON gm.group_id = gc.group_id
						WHERE gc.course_id = courses.id
						  AND gm.user_id = ?
					)
				)
			)
		`, userId).
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

	_ = r.redis.Del(ctx,
		"course:"+course.Id,
		"courses:all",
	).Err()

	return &updatedCourse, nil
}

func (r *courseRepository) UpdatePublishedAt(ctx context.Context, id string, time *time.Time) error {
	err := r.db.
		WithContext(ctx).
		Model(&entities.Course{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"published_at": time,
		}).Error

	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"course:"+id,
		"courses:all",
	).Err()

	return nil
}

func (r *courseRepository) Delete(ctx context.Context, id string) error {
	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&entities.Course{}).Error

	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"course:"+id,
		"courses:all",
	).Err()

	return nil
}
