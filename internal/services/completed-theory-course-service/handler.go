package completed_theory_course_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/completedtheorycourseservicepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompletedTheoryCourseService interface {
	CreateCompletedTheoryCourse(ctx context.Context, completedCourse *entities.CompletedTheoryCourse) error
	ReadCompletedTheoryCourseByUserIdAndCourseId(ctx context.Context, userId string, lessonId string) (*entities.CompletedTheoryCourse, error)
	ReadAllCompletedTheoryCoursesByCourseId(ctx context.Context, courseId string) ([]entities.CompletedTheoryCourse, error)
	ReadAllCompletedTheoryCoursesByUserId(ctx context.Context, userId string) ([]entities.CompletedTheoryCourse, error)
}

type CompletedTheoryCourseHandler struct {
	completedtheorycourseservicepb.UnimplementedCompletedTheoryCourseServiceServer
	service CompletedTheoryCourseService
}

func NewCompletedTheoryCourseHandler(service CompletedTheoryCourseService) *CompletedTheoryCourseHandler {
	return &CompletedTheoryCourseHandler{service: service}
}

func (h *CompletedTheoryCourseHandler) CreateCompletedTheoryCourse(ctx context.Context, req *completedtheorycourseservicepb.CreateCompletedTheoryCourseRequest) (*completedtheorycourseservicepb.CreateCompletedTheoryCourseResponse, error) {
	completedCourse := entities.CompletedTheoryCourse{
		UserId:   req.UserId,
		CourseId: req.CourseId,
	}

	err := h.service.CreateCompletedTheoryCourse(ctx, &completedCourse)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedtheorycourseservicepb.CreateCompletedTheoryCourseResponse{
		CompletedTheoryCourse: &completedtheorycourseservicepb.CompletedTheoryCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
	}, nil
}

func (h *CompletedTheoryCourseHandler) ReadCompletedTheoryCourseByUserIdAndCourseId(ctx context.Context, req *completedtheorycourseservicepb.ReadCompletedTheoryCourseByUserIdAndCourseIdRequest) (*completedtheorycourseservicepb.ReadCompletedTheoryCourseByUserIdAndCourseIdResponse, error) {
	completedCourse, err := h.service.ReadCompletedTheoryCourseByUserIdAndCourseId(ctx, req.UserId, req.CourseId)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return &completedtheorycourseservicepb.ReadCompletedTheoryCourseByUserIdAndCourseIdResponse{}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedtheorycourseservicepb.ReadCompletedTheoryCourseByUserIdAndCourseIdResponse{
		CompletedTheoryCourse: &completedtheorycourseservicepb.CompletedTheoryCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
	}, nil
}

func (h *CompletedTheoryCourseHandler) ReadAllCompletedTheoryCoursesByCourseId(ctx context.Context, req *completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByCourseIdRequest) (*completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByCourseIdResponse, error) {
	completedTheoryCourses, err := h.service.ReadAllCompletedTheoryCoursesByCourseId(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedTheoryCoursesPb := make([]*completedtheorycourseservicepb.CompletedTheoryCourse, 0, len(completedTheoryCourses))
	for _, completedCourse := range completedTheoryCourses {
		completedTheoryCoursesPb = append(completedTheoryCoursesPb, &completedtheorycourseservicepb.CompletedTheoryCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
		)
	}

	return &completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByCourseIdResponse{
		CompletedTheoryCourses: completedTheoryCoursesPb,
	}, nil
}

func (h *CompletedTheoryCourseHandler) ReadAllCompletedTheoryCoursesByUserId(ctx context.Context, req *completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByUserIdRequest) (*completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByUserIdResponse, error) {
	completedTheoryCourses, err := h.service.ReadAllCompletedTheoryCoursesByUserId(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedTheoryCoursesPb := make([]*completedtheorycourseservicepb.CompletedTheoryCourse, 0, len(completedTheoryCourses))
	for _, completedCourse := range completedTheoryCourses {
		completedTheoryCoursesPb = append(completedTheoryCoursesPb, &completedtheorycourseservicepb.CompletedTheoryCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
		)
	}

	return &completedtheorycourseservicepb.ReadAllCompletedTheoryCoursesByUserIdResponse{
		CompletedTheoryCourses: completedTheoryCoursesPb,
	}, nil
}
