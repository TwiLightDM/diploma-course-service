package task_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *taskRepository) Create(ctx context.Context, task *entities.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errs.ErrDuplicateKey
		}

		return err
	}

	return nil
}

func (r *taskRepository) ReadById(ctx context.Context, id string) (*entities.Task, error) {
	var task entities.Task

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&task)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.ErrRecordNotFound
		}

		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) ReadAllByCourseId(ctx context.Context, courseId string) ([]entities.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"course_id": courseId,
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tasks []entities.Task

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) ReadAllByModuleId(ctx context.Context, moduleId string) ([]entities.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"module_id": moduleId,
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tasks []entities.Task

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	set := bson.M{}

	if task.Title != "" {
		set["title"] = task.Title
	}

	if task.CorrectTextAnswer != "" {
		set["correct_text_answer"] = task.CorrectTextAnswer
	}

	if task.Options != nil {
		set["options"] = task.Options
	}

	if len(set) == 0 {
		return r.ReadById(ctx, task.Id)
	}

	update := bson.M{
		"$set": set,
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id": task.Id,
	}, update)

	if err != nil {
		return nil, err
	}

	return r.ReadById(ctx, task.Id)
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errs.ErrRecordNotFound
	}

	return nil
}
