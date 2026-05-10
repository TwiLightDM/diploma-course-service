package task_attempt_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type taskAttemptRepository struct {
	collection *mongo.Collection
}

func NewTaskAttemptRepository(db *mongo.Database) TaskAttemptRepository {
	return &taskAttemptRepository{collection: db.Collection("task_attempts")}
}

func (r *taskAttemptRepository) Create(ctx context.Context, taskAttempt *entities.TaskAttempt) error {
	_, err := r.collection.InsertOne(ctx, taskAttempt)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errs.ErrDuplicateKey
		}

		return err
	}

	return nil
}

func (r *taskAttemptRepository) ReadById(ctx context.Context, id string) (*entities.TaskAttempt, error) {
	var taskAttempt entities.TaskAttempt

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&taskAttempt)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	return &taskAttempt, nil
}

func (r *taskAttemptRepository) ReadAllByUserIdAndModuleId(ctx context.Context, userId, moduleId string) ([]entities.TaskAttempt, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"user_id":   userId,
			"module_id": moduleId,
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var taskAttempts []entities.TaskAttempt

	if err = cursor.All(ctx, &taskAttempts); err != nil {
		return nil, err
	}

	return taskAttempts, nil
}

func (r *taskAttemptRepository) ReadAllByUserIdAndCourseId(ctx context.Context, userId, courseId string) ([]entities.TaskAttempt, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"user_id":   userId,
			"course_id": courseId,
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var taskAttempts []entities.TaskAttempt

	if err = cursor.All(ctx, &taskAttempts); err != nil {
		return nil, err
	}

	return taskAttempts, nil
}

func (r *taskAttemptRepository) GetLastAttemptNumber(ctx context.Context, userId string, courseId string, moduleId string) (int, error) {
	filter := bson.M{
		"user_id": userId,
	}

	if courseId != "" {
		filter["course_id"] = courseId
	}

	if moduleId != "" {
		filter["module_id"] = moduleId
	}

	opts := options.FindOne().
		SetSort(bson.D{
			{
				Key:   "attempt_number",
				Value: -1,
			},
		})

	var taskAttempt entities.TaskAttempt

	err := r.collection.FindOne(ctx, filter, opts).Decode(&taskAttempt)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}

		return 0, err
	}

	return taskAttempt.AttemptNumber, nil
}
