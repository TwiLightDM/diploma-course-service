package group_course_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"gorm.io/gorm"
)

type GroupCourseRepository interface {
	Create(ctx context.Context, groupCourse *entities.GroupCourse) error
	ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.GroupCourse, error)
	ReadAllByGroupId(ctx context.Context, groupId string) ([]entities.GroupCourse, error)
	Delete(ctx context.Context, id string) error
}

type groupCourseRepository struct {
	db *gorm.DB
}

func NewGroupCourseRepository(db *gorm.DB) GroupCourseRepository {
	return &groupCourseRepository{db: db}
}

func (r *groupCourseRepository) Create(ctx context.Context, groupCourse *entities.GroupCourse) error {
	return r.db.WithContext(ctx).Create(groupCourse).Error
}

func (r *groupCourseRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.GroupCourse, error) {
	var groupCourses []entities.GroupCourse
	if err := r.db.
		WithContext(ctx).
		Where("course_id = ?", courseId).
		Find(&groupCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return groupCourses, nil
}

func (r *groupCourseRepository) ReadAllByGroupId(ctx context.Context, groupId string) ([]entities.GroupCourse, error) {
	var groupCourses []entities.GroupCourse
	if err := r.db.
		WithContext(ctx).
		Where("group_id = ?", groupId).
		Find(&groupCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return groupCourses, nil
}

func (r *groupCourseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.GroupCourse{}).Error
}
