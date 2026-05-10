package task_attempt_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/taskattemptservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskAttemptService interface {
	SubmitTaskAttempt(ctx context.Context, taskAttempt *entities.TaskAttempt) (*entities.TaskAttempt, error)
	ReadTaskAttemptById(ctx context.Context, id string) (*entities.TaskAttempt, error)
	ReadAllTaskAttemptsByUserIdAndModuleId(ctx context.Context, userId, moduleId string) ([]entities.TaskAttempt, error)
	ReadAllTaskAttemptsByUserIdAndCourseId(ctx context.Context, userId, courseId string) ([]entities.TaskAttempt, error)
}

type TaskAttemptHandler struct {
	taskattemptservicepb.UnimplementedTaskAttemptServiceServer
	service TaskAttemptService
}

func NewTaskAttemptHandler(service TaskAttemptService) *TaskAttemptHandler {
	return &TaskAttemptHandler{service: service}
}

func (h *TaskAttemptHandler) SubmitTaskAttempt(ctx context.Context, req *taskattemptservicepb.SubmitTaskAttemptRequest) (*taskattemptservicepb.SubmitTaskAttemptResponse, error) {
	taskAttempt := &entities.TaskAttempt{
		UserId:   req.UserId,
		CourseId: req.CourseId,
		ModuleId: req.ModuleId,
		Answers:  make([]entities.TaskAttemptAnswer, 0, len(req.Answers)),
	}

	if (taskAttempt.CourseId == "" && taskAttempt.ModuleId == "") || (taskAttempt.CourseId != "" && taskAttempt.ModuleId != "") {
		return nil, status.Error(codes.InvalidArgument, "module id or course id is required or they can not be both filled")
	}

	for _, answer := range req.Answers {
		taskAttempt.Answers = append(taskAttempt.Answers, entities.TaskAttemptAnswer{
			TaskId:            answer.TaskId,
			TextAnswer:        answer.TextAnswer,
			SelectedOptionIds: answer.SelectedOptionIds,
		})
	}

	createdTaskAttempt, err := h.service.SubmitTaskAttempt(ctx, taskAttempt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskattemptservicepb.SubmitTaskAttemptResponse{
		TaskAttempt: mapTaskAttemptToProto(createdTaskAttempt),
	}, nil
}

func (h *TaskAttemptHandler) ReadTaskAttempt(ctx context.Context, req *taskattemptservicepb.ReadTaskAttemptRequest) (*taskattemptservicepb.ReadTaskAttemptResponse, error) {
	taskAttempt, err := h.service.ReadTaskAttemptById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskattemptservicepb.ReadTaskAttemptResponse{
		TaskAttempt: mapTaskAttemptToProto(taskAttempt),
	}, nil
}

func (h *TaskAttemptHandler) ReadAllTaskAttemptsByUserIdAndModuleId(ctx context.Context, req *taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndModuleIdRequest) (*taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndModuleIdResponse, error) {
	taskAttempts, err := h.service.ReadAllTaskAttemptsByUserIdAndModuleId(ctx, req.UserId, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	taskAttemptsPb := make([]*taskattemptservicepb.TaskAttempt, 0, len(taskAttempts))
	for _, taskAttempt := range taskAttempts {
		taskAttemptsPb = append(taskAttemptsPb, mapTaskAttemptToProto(&taskAttempt))
	}

	return &taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndModuleIdResponse{
		TaskAttempts: taskAttemptsPb,
	}, nil
}

func (h *TaskAttemptHandler) ReadAllTaskAttemptsByUserIdAndCourseId(ctx context.Context, req *taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndCourseIdRequest) (*taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndCourseIdResponse, error) {
	taskAttempts, err := h.service.ReadAllTaskAttemptsByUserIdAndCourseId(ctx, req.UserId, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	taskAttemptsPb := make([]*taskattemptservicepb.TaskAttempt, 0, len(taskAttempts))
	for _, taskAttempt := range taskAttempts {
		taskAttemptsPb = append(taskAttemptsPb, mapTaskAttemptToProto(&taskAttempt))
	}

	return &taskattemptservicepb.ReadAllTaskAttemptsByUserIdAndCourseIdResponse{
		TaskAttempts: taskAttemptsPb,
	}, nil
}

func mapTaskAttemptToProto(taskAttempt *entities.TaskAttempt) *taskattemptservicepb.TaskAttempt {
	answers := make([]*taskattemptservicepb.TaskAttemptAnswer, 0, len(taskAttempt.Answers))
	for _, answer := range taskAttempt.Answers {
		answerPb := &taskattemptservicepb.TaskAttemptAnswer{
			TaskId:            answer.TaskId,
			SelectedOptionIds: answer.SelectedOptionIds,
			TextAnswer:        answer.TextAnswer,
			IsCorrect:         answer.IsCorrect,
		}

		answers = append(answers, answerPb)
	}

	taskAttemptPb := &taskattemptservicepb.TaskAttempt{
		Id:             taskAttempt.Id,
		UserId:         taskAttempt.UserId,
		CourseId:       taskAttempt.CourseId,
		ModuleId:       taskAttempt.ModuleId,
		AttemptNumber:  int64(taskAttempt.AttemptNumber),
		Answers:        answers,
		CorrectAnswers: int64(taskAttempt.CorrectAnswers),
		TotalQuestions: int64(taskAttempt.TotalQuestions),
		Score:          taskAttempt.Score,
	}

	return taskAttemptPb
}
