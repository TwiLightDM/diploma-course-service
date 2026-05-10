package task_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/taskservicepb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *entities.Task) error
	ReadTaskById(ctx context.Context, id string) (*entities.Task, error)
	ReadAllTasksByCourseId(ctx context.Context, courseId string) ([]entities.Task, error)
	ReadAllTasksByModuleId(ctx context.Context, moduleId string) ([]entities.Task, error)
	UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type TaskHandler struct {
	taskservicepb.UnimplementedTaskServiceServer
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *taskservicepb.CreateTaskRequest) (*taskservicepb.CreateTaskResponse, error) {
	task := &entities.Task{
		Title:    req.Title,
		CourseId: req.CourseId,
		ModuleId: req.ModuleId,
		Options:  make([]entities.TaskOption, 0, len(req.Options)),
	}

	if (task.CourseId == "" && task.ModuleId == "") || (task.CourseId != "" && task.ModuleId != "") {
		return nil, status.Error(codes.InvalidArgument, "module id or course id is required or they can not be both filled")
	}

	switch req.Type {
	case taskservicepb.TaskType_TASK_TYPE_TEXT_INPUT:
		task.Type = entities.TaskTypeTextInput

		task.CorrectTextAnswer = req.CorrectTextAnswer

	case taskservicepb.TaskType_TASK_TYPE_CHOICE:
		task.Type = entities.TaskTypeChoice

		for _, option := range req.Options {
			task.Options = append(task.Options, entities.TaskOption{
				Id:        uuid.NewString(),
				Text:      option.Text,
				IsCorrect: option.IsCorrect,
			})
		}
	}

	err := h.service.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservicepb.CreateTaskResponse{
		Task: mapTaskToProto(task),
	}, nil
}

func (h *TaskHandler) ReadTask(ctx context.Context, req *taskservicepb.ReadTaskRequest) (*taskservicepb.ReadTaskResponse, error) {
	task, err := h.service.ReadTaskById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservicepb.ReadTaskResponse{
		Task: mapTaskToProto(task),
	}, nil
}

func (h *TaskHandler) ReadAllTasksByCourseId(ctx context.Context, req *taskservicepb.ReadAllTasksByCourseIdRequest) (*taskservicepb.ReadAllTasksByCourseIdResponse, error) {
	tasks, err := h.service.ReadAllTasksByCourseId(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tasksPb := make([]*taskservicepb.Task, 0, len(tasks))

	for _, task := range tasks {
		tasksPb = append(tasksPb, mapTaskToProto(&task))
	}

	return &taskservicepb.ReadAllTasksByCourseIdResponse{
		Tasks: tasksPb,
	}, nil
}

func (h *TaskHandler) ReadAllTasksByModuleId(ctx context.Context, req *taskservicepb.ReadAllTasksByModuleIdRequest) (*taskservicepb.ReadAllTasksByModuleIdResponse, error) {
	tasks, err := h.service.ReadAllTasksByModuleId(ctx, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tasksPb := make([]*taskservicepb.Task, 0, len(tasks))

	for _, task := range tasks {
		tasksPb = append(tasksPb, mapTaskToProto(&task))
	}

	return &taskservicepb.ReadAllTasksByModuleIdResponse{
		Tasks: tasksPb,
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *taskservicepb.UpdateTaskRequest) (*taskservicepb.UpdateTaskResponse, error) {
	task := &entities.Task{
		Id:                req.Id,
		Title:             req.Title,
		CorrectTextAnswer: req.CorrectTextAnswer,
		Options:           make([]entities.TaskOption, 0, len(req.Options)),
	}

	for _, option := range req.Options {
		task.Options = append(task.Options, entities.TaskOption{
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
		})
	}

	updatedTask, err := h.service.UpdateTask(ctx, task)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservicepb.UpdateTaskResponse{
		Task: mapTaskToProto(updatedTask),
	}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *taskservicepb.DeleteTaskRequest) (*taskservicepb.DeleteTaskResponse, error) {
	if err := h.service.DeleteTask(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservicepb.DeleteTaskResponse{}, nil
}

func mapTaskToProto(task *entities.Task) *taskservicepb.Task {
	taskPb := &taskservicepb.Task{
		Id:       task.Id,
		Title:    task.Title,
		CourseId: task.CourseId,
		ModuleId: task.ModuleId,
	}

	switch task.Type {
	case entities.TaskTypeTextInput:
		taskPb.Type = taskservicepb.TaskType_TASK_TYPE_TEXT_INPUT
		taskPb.CorrectTextAnswer = task.CorrectTextAnswer

	case entities.TaskTypeChoice:
		taskPb.Type = taskservicepb.TaskType_TASK_TYPE_CHOICE
		options := make([]*taskservicepb.TaskOption, 0, len(task.Options))

		for _, option := range task.Options {
			options = append(options, &taskservicepb.TaskOption{
				Id:        option.Id,
				Text:      option.Text,
				IsCorrect: option.IsCorrect,
			})
		}

		taskPb.Options = options
	}

	return taskPb
}
