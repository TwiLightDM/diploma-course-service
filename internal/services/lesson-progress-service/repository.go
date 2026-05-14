package lesson_progress_service

import (
	"context"
	"errors"
	"strings"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"

	"gorm.io/gorm"
)

type lessonProgressRepository struct {
	db *gorm.DB
}

func NewLessonProgressRepository(db *gorm.DB) LessonProgressRepository {
	return &lessonProgressRepository{db: db}
}

func (r *lessonProgressRepository) Create(ctx context.Context, progress *entities.LessonProgress) error {
	err := r.db.WithContext(ctx).Create(progress).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") || strings.Contains(err.Error(), "SQLSTATE 23505") {
			return errs.ErrDuplicateKey
		}

		return err
	}

	return nil
}

func (r *lessonProgressRepository) ReadByUserId(ctx context.Context, userId string) ([]entities.LessonProgress, error) {
	var progresses []entities.LessonProgress

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&progresses).Error

	if err != nil {
		return nil, err
	}

	return progresses, nil
}

func (r *lessonProgressRepository) ReadByUserIdAndLessonId(ctx context.Context, userId string, lessonId string) (*entities.LessonProgress, error) {
	var progress entities.LessonProgress

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND lesson_id = ?", userId, lessonId).
		First(&progress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	return &progress, nil
}

func (r *lessonProgressRepository) ReadCourseProgressByUserId(ctx context.Context, userId string, courseId string) (*entities.CourseProgress, error) {
	type result struct {
		TotalLessons     int64
		CompletedLessons int64
	}

	var res result

	err := r.db.WithContext(ctx).
		Table("lessons").
		Select(`
			COUNT(lessons.id) AS total_lessons,
			COUNT(lp.lesson_id) AS completed_lessons
		`).
		Joins(`
			JOIN modules 
				ON modules.id = lessons.module_id
		`).
		Joins(`
			LEFT JOIN lesson_progresses lp
				ON lp.lesson_id = lessons.id
				AND lp.user_id = ?
		`, userId).
		Where(`
			modules.course_id = ?
			AND modules.deleted_at IS NULL
			AND lessons.deleted_at IS NULL
		`, courseId).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	var completedLessonIds []string

	err = r.db.WithContext(ctx).
		Table("lesson_progresses lp").
		Select("lp.lesson_id").
		Joins(`
			JOIN lessons 
				ON lessons.id = lp.lesson_id
		`).
		Joins(`
			JOIN modules
				ON modules.id = lessons.module_id
		`).
		Where(`
			lp.user_id = ?
			AND modules.course_id = ?
		`, userId, courseId).
		Pluck("lp.lesson_id", &completedLessonIds).Error

	if err != nil {
		return nil, err
	}

	var progressPercent float64

	if res.TotalLessons > 0 {
		progressPercent = float64(res.CompletedLessons) / float64(res.TotalLessons) * 100
	}

	return &entities.CourseProgress{
		CourseId:           courseId,
		TotalLessons:       res.TotalLessons,
		CompletedLessons:   res.CompletedLessons,
		ProgressPercent:    progressPercent,
		CompletedLessonIds: completedLessonIds,
	}, nil
}

func (r *lessonProgressRepository) ReadCourseStatistics(ctx context.Context, courseId string) ([]entities.UserCourseProgress, error) {
	type result struct {
		UserId           string
		CompletedLessons int64
		TotalLessons     int64
		ProgressPercent  float64
	}

	var results []result

	totalLessonsSubQuery := r.db.
		Table("lessons").
		Select("COUNT(lessons.id)").
		Joins(`
			JOIN modules
				ON modules.id = lessons.module_id
		`).
		Where(`
			modules.course_id = ?
			AND modules.deleted_at IS NULL
			AND lessons.deleted_at IS NULL
		`, courseId)

	err := r.db.WithContext(ctx).
		Table("lesson_progresses lp").
		Select(`
			lp.user_id,
			COUNT(lp.lesson_id) AS completed_lessons,
			(?) AS total_lessons,
			COUNT(lp.lesson_id) * 100.0 / (?) AS progress_percent
		`, totalLessonsSubQuery, totalLessonsSubQuery).
		Joins(`
			JOIN lessons
				ON lessons.id = lp.lesson_id
		`).
		Joins(`
			JOIN modules
				ON modules.id = lessons.module_id
		`).
		Where(`
			modules.course_id = ?
		`, courseId).
		Group("lp.user_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	progresses := make([]entities.UserCourseProgress, 0, len(results))
	for _, res := range results {

		progresses = append(progresses, entities.UserCourseProgress{
			UserId:           res.UserId,
			CompletedLessons: res.CompletedLessons,
			TotalLessons:     res.TotalLessons,
			ProgressPercent:  res.ProgressPercent,
			Completed:        res.TotalLessons > 0 && res.CompletedLessons == res.TotalLessons,
		},
		)
	}

	return progresses, nil
}

func (r *lessonProgressRepository) ReadModuleProgressByUserId(ctx context.Context, userId string, moduleId string) (*entities.ModuleProgress, error) {
	type result struct {
		TotalLessons     int64
		CompletedLessons int64
	}

	var res result

	err := r.db.WithContext(ctx).
		Table("lessons").
		Select(`
			COUNT(lessons.id) AS total_lessons,
			COUNT(lp.lesson_id) AS completed_lessons
		`).
		Joins(`
			LEFT JOIN lesson_progresses lp
				ON lp.lesson_id = lessons.id
				AND lp.user_id = ?
		`, userId).
		Where(`
			lessons.module_id = ?
			AND lessons.deleted_at IS NULL
		`, moduleId).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	var completedLessonIds []string

	err = r.db.WithContext(ctx).
		Table("lesson_progresses").
		Where(`
			user_id = ?
			AND lesson_id IN (
				SELECT id
				FROM lessons
				WHERE module_id = ?
			)
		`, userId, moduleId).
		Pluck("lesson_id", &completedLessonIds).Error

	if err != nil {
		return nil, err
	}

	var progressPercent float64

	if res.TotalLessons > 0 {
		progressPercent = float64(res.CompletedLessons) / float64(res.TotalLessons) * 100
	}

	return &entities.ModuleProgress{
		ModuleId:           moduleId,
		TotalLessons:       res.TotalLessons,
		CompletedLessons:   res.CompletedLessons,
		ProgressPercent:    progressPercent,
		CompletedLessonIds: completedLessonIds,
	}, nil
}

func (r *lessonProgressRepository) ReadModuleStatistics(ctx context.Context, moduleId string) ([]entities.UserModuleProgress, error) {
	type result struct {
		UserId           string
		CompletedLessons int64
		TotalLessons     int64
		ProgressPercent  float64
	}

	var results []result

	totalLessonsSubQuery := r.db.
		Table("lessons").
		Select("COUNT(id)").
		Where(`
			module_id = ?
			AND deleted_at IS NULL
		`, moduleId)

	err := r.db.WithContext(ctx).
		Table("lesson_progresses lp").
		Select(`
			lp.user_id,
			COUNT(lp.lesson_id) AS completed_lessons,
			(?) AS total_lessons,
			COUNT(lp.lesson_id) * 100.0 / (?) AS progress_percent
		`, totalLessonsSubQuery, totalLessonsSubQuery).
		Joins(`
			JOIN lessons
				ON lessons.id = lp.lesson_id
		`).
		Where(`
			lessons.module_id = ?
		`, moduleId).
		Group("lp.user_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	progresses := make([]entities.UserModuleProgress, 0, len(results))
	for _, res := range results {
		progresses = append(progresses, entities.UserModuleProgress{
			UserId:           res.UserId,
			CompletedLessons: res.CompletedLessons,
			TotalLessons:     res.TotalLessons,
			ProgressPercent:  res.ProgressPercent,
			Completed:        res.TotalLessons > 0 && res.CompletedLessons == res.TotalLessons,
		},
		)
	}

	return progresses, nil
}
