package completed_course_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/completedcourseservicepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompletedCourseService interface {
	CreateCompletedCourse(ctx context.Context, completedCourse *entities.CompletedCourse) error
	ReadCompletedCourseByUserIdAndCourseId(ctx context.Context, userId string, lessonId string) (*entities.CompletedCourse, error)
	ReadAllCompletedCoursesByCourseId(ctx context.Context, courseId string) ([]entities.CompletedCourse, error)
	ReadAllCompletedCoursesByUserId(ctx context.Context, userId string) ([]entities.CompletedCourse, error)
}

type CompletedCourseHandler struct {
	completedcourseservicepb.UnimplementedCompletedCourseServiceServer
	service CompletedCourseService
}

func NewCompletedCourseHandler(service CompletedCourseService) *CompletedCourseHandler {
	return &CompletedCourseHandler{service: service}
}

func (h *CompletedCourseHandler) CreateCompletedCourse(ctx context.Context, req *completedcourseservicepb.CreateCompletedCourseRequest) (*completedcourseservicepb.CreateCompletedCourseResponse, error) {
	completedCourse := entities.CompletedCourse{
		UserId:   req.UserId,
		CourseId: req.CourseId,
	}

	err := h.service.CreateCompletedCourse(ctx, &completedCourse)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedcourseservicepb.CreateCompletedCourseResponse{
		CompletedCourse: &completedcourseservicepb.CompletedCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
	}, nil
}

func (h *CompletedCourseHandler) ReadCompletedCourseByUserIdAndCourseId(ctx context.Context, req *completedcourseservicepb.ReadCompletedCourseByUserIdAndCourseIdRequest) (*completedcourseservicepb.ReadCompletedCourseByUserIdAndCourseIdResponse, error) {
	completedCourse, err := h.service.ReadCompletedCourseByUserIdAndCourseId(ctx, req.UserId, req.CourseId)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return &completedcourseservicepb.ReadCompletedCourseByUserIdAndCourseIdResponse{}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedcourseservicepb.ReadCompletedCourseByUserIdAndCourseIdResponse{
		CompletedCourse: &completedcourseservicepb.CompletedCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
	}, nil
}

func (h *CompletedCourseHandler) ReadAllCompletedCoursesByCourseId(ctx context.Context, req *completedcourseservicepb.ReadAllCompletedCoursesByCourseIdRequest) (*completedcourseservicepb.ReadAllCompletedCoursesByCourseIdResponse, error) {
	completedCourses, err := h.service.ReadAllCompletedCoursesByCourseId(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedCoursesPb := make([]*completedcourseservicepb.CompletedCourse, 0, len(completedCourses))
	for _, completedCourse := range completedCourses {
		completedCoursesPb = append(completedCoursesPb, &completedcourseservicepb.CompletedCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
		)
	}

	return &completedcourseservicepb.ReadAllCompletedCoursesByCourseIdResponse{
		CompletedCourses: completedCoursesPb,
	}, nil
}

func (h *CompletedCourseHandler) ReadAllCompletedCoursesByUserId(ctx context.Context, req *completedcourseservicepb.ReadAllCompletedCoursesByUserIdRequest) (*completedcourseservicepb.ReadAllCompletedCoursesByUserIdResponse, error) {
	completedCourses, err := h.service.ReadAllCompletedCoursesByUserId(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedCoursesPb := make([]*completedcourseservicepb.CompletedCourse, 0, len(completedCourses))
	for _, completedCourse := range completedCourses {
		completedCoursesPb = append(completedCoursesPb, &completedcourseservicepb.CompletedCourse{
			UserId:   completedCourse.UserId,
			CourseId: completedCourse.CourseId,
		},
		)
	}

	return &completedcourseservicepb.ReadAllCompletedCoursesByUserIdResponse{
		CompletedCourses: completedCoursesPb,
	}, nil
}
