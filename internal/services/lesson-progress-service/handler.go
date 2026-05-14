package lesson_progress_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonprogressservicepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LessonProgressService interface {
	CreateLessonProgress(ctx context.Context, progress *entities.LessonProgress) error
	ReadLessonProgressByUserId(ctx context.Context, userId string) ([]entities.LessonProgress, error)
	ReadLessonProgressByUserIdAndLessonId(ctx context.Context, userId string, lessonId string) (*entities.LessonProgress, error)
	ReadModuleProgressByUserId(ctx context.Context, userId string, moduleId string) (*entities.ModuleProgress, error)
	ReadModuleStatistics(ctx context.Context, moduleId string) ([]entities.UserModuleProgress, error)
	ReadCourseProgressByUserId(ctx context.Context, userId string, courseId string) (*entities.CourseProgress, error)
	ReadCourseStatistics(ctx context.Context, courseId string) ([]entities.UserCourseProgress, error)
}

type LessonProgressHandler struct {
	lessonprogressservicepb.UnimplementedLessonProgressServiceServer
	service LessonProgressService
}

func NewLessonProgressHandler(service LessonProgressService) *LessonProgressHandler {
	return &LessonProgressHandler{service: service}
}

func (h *LessonProgressHandler) CreateLessonProgress(ctx context.Context, req *lessonprogressservicepb.CreateLessonProgressRequest) (*lessonprogressservicepb.CreateLessonProgressResponse, error) {
	progress := entities.LessonProgress{
		UserId:   req.UserId,
		LessonId: req.LessonId,
	}

	err := h.service.CreateLessonProgress(ctx, &progress)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonprogressservicepb.CreateLessonProgressResponse{
		LessonProgress: &lessonprogressservicepb.LessonProgress{
			UserId:   progress.UserId,
			LessonId: progress.LessonId,
		},
	}, nil
}

func (h *LessonProgressHandler) ReadLessonProgressByUserId(ctx context.Context, req *lessonprogressservicepb.ReadLessonProgressByUserIdRequest) (*lessonprogressservicepb.ReadLessonProgressByUserIdResponse, error) {
	progresses, err := h.service.ReadLessonProgressByUserId(ctx, req.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	progressesPb := make([]*lessonprogressservicepb.LessonProgress, 0, len(progresses))

	for _, progress := range progresses {
		progressesPb = append(progressesPb, &lessonprogressservicepb.LessonProgress{
			UserId:   progress.UserId,
			LessonId: progress.LessonId,
		},
		)
	}

	return &lessonprogressservicepb.ReadLessonProgressByUserIdResponse{
		LessonProgresses: progressesPb,
	}, nil
}

func (h *LessonProgressHandler) ReadLessonProgressByUserIdAndLessonId(ctx context.Context, req *lessonprogressservicepb.ReadLessonProgressByUserIdAndLessonIdRequest) (*lessonprogressservicepb.ReadLessonProgressByUserIdAndLessonIdResponse, error) {
	progress, err := h.service.ReadLessonProgressByUserIdAndLessonId(ctx, req.UserId, req.LessonId)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return &lessonprogressservicepb.ReadLessonProgressByUserIdAndLessonIdResponse{}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonprogressservicepb.ReadLessonProgressByUserIdAndLessonIdResponse{
		LessonProgress: &lessonprogressservicepb.LessonProgress{
			UserId:   progress.UserId,
			LessonId: progress.LessonId,
		},
	}, nil
}

func (h *LessonProgressHandler) ReadCourseProgressByUserId(ctx context.Context, req *lessonprogressservicepb.ReadCourseProgressByUserIdRequest) (*lessonprogressservicepb.ReadCourseProgressByUserIdResponse, error) {
	progress, err := h.service.ReadCourseProgressByUserId(ctx, req.UserId, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonprogressservicepb.ReadCourseProgressByUserIdResponse{
		CourseProgress: &lessonprogressservicepb.CourseProgress{
			CourseId:           progress.CourseId,
			TotalLessons:       progress.TotalLessons,
			CompletedLessons:   progress.CompletedLessons,
			ProgressPercent:    progress.ProgressPercent,
			CompletedLessonIds: progress.CompletedLessonIds,
		},
	}, nil
}

func (h *LessonProgressHandler) ReadCourseStatistics(ctx context.Context, req *lessonprogressservicepb.ReadCourseStatisticsRequest) (*lessonprogressservicepb.ReadCourseStatisticsResponse, error) {
	usersProgress, err := h.service.ReadCourseStatistics(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	usersProgressPb := make([]*lessonprogressservicepb.UserCourseProgress, 0, len(usersProgress))
	for _, progress := range usersProgress {
		usersProgressPb = append(usersProgressPb, &lessonprogressservicepb.UserCourseProgress{
			UserId:           progress.UserId,
			CompletedLessons: progress.CompletedLessons,
			TotalLessons:     progress.TotalLessons,
			ProgressPercent:  progress.ProgressPercent,
			Completed:        progress.Completed,
		},
		)
	}

	return &lessonprogressservicepb.ReadCourseStatisticsResponse{
		CourseStatistics: &lessonprogressservicepb.CourseStatistics{
			CourseId:      req.CourseId,
			UsersProgress: usersProgressPb,
		},
	}, nil
}

func (h *LessonProgressHandler) ReadModuleProgressByUserId(ctx context.Context, req *lessonprogressservicepb.ReadModuleProgressByUserIdRequest) (*lessonprogressservicepb.ReadModuleProgressByUserIdResponse, error) {
	progress, err := h.service.ReadModuleProgressByUserId(ctx, req.UserId, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonprogressservicepb.ReadModuleProgressByUserIdResponse{
		ModuleProgress: &lessonprogressservicepb.ModuleProgress{
			ModuleId:           progress.ModuleId,
			TotalLessons:       progress.TotalLessons,
			CompletedLessons:   progress.CompletedLessons,
			ProgressPercent:    progress.ProgressPercent,
			CompletedLessonIds: progress.CompletedLessonIds,
		},
	}, nil
}

func (h *LessonProgressHandler) ReadModuleStatistics(ctx context.Context, req *lessonprogressservicepb.ReadModuleStatisticsRequest) (*lessonprogressservicepb.ReadModuleStatisticsResponse, error) {
	usersProgress, err := h.service.ReadModuleStatistics(ctx, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	usersProgressPb := make([]*lessonprogressservicepb.UserModuleProgress, 0, len(usersProgress))
	for _, progress := range usersProgress {
		usersProgressPb = append(usersProgressPb, &lessonprogressservicepb.UserModuleProgress{
			UserId:           progress.UserId,
			CompletedLessons: progress.CompletedLessons,
			TotalLessons:     progress.TotalLessons,
			ProgressPercent:  progress.ProgressPercent,
			Completed:        progress.Completed,
		},
		)
	}

	return &lessonprogressservicepb.ReadModuleStatisticsResponse{
		ModuleStatistics: &lessonprogressservicepb.ModuleStatistics{
			ModuleId:      req.ModuleId,
			UsersProgress: usersProgressPb,
		},
	}, nil
}
