package course_service

import (
	"context"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/courseservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourseHandler struct {
	courseservicepb.UnimplementedCourseServiceServer
	service CourseService
}

func NewCourseHandler(service CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) CreateCourse(ctx context.Context, req *courseservicepb.CreateCourseRequest) (*courseservicepb.CreateCourseResponse, error) {
	course := entities.Course{
		Title:       req.Title,
		Description: req.Description,
		AccessType:  req.AccessType,
		OwnerId:     req.OwnerId,
	}

	err := h.service.CreateCourse(ctx, &course)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &courseservicepb.CreateCourseResponse{
		Course: &courseservicepb.Course{
			Id:          course.Id,
			Title:       course.Title,
			Description: course.Description,
			AccessType:  course.AccessType,
			IsPublished: course.IsPublished,
			OwnerId:     course.OwnerId,
		},
	}, nil
}

func (h *CourseHandler) ReadCourse(ctx context.Context, req *courseservicepb.ReadCourseRequest) (*courseservicepb.ReadCourseResponse, error) {
	course, err := h.service.ReadCourseById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &courseservicepb.ReadCourseResponse{
		Course: &courseservicepb.Course{
			Id:          course.Id,
			Title:       course.Title,
			Description: course.Description,
			AccessType:  course.AccessType,
			OwnerId:     course.OwnerId,
			IsPublished: course.IsPublished,
		},
	}, nil
}

func (h *CourseHandler) ReadAllCoursesByOwnerId(ctx context.Context, req *courseservicepb.ReadAllCoursesByOwnerIdRequest) (*courseservicepb.ReadAllCoursesByOwnerIdResponse, error) {
	courses, err := h.service.ReadAllCoursesByOwnerId(ctx, req.OwnerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	coursesPb := make([]*courseservicepb.Course, 0, len(courses))
	for _, course := range courses {
		coursesPb = append(coursesPb, &courseservicepb.Course{
			Id:          course.Id,
			Title:       course.Title,
			Description: course.Description,
			AccessType:  course.AccessType,
			IsPublished: course.IsPublished,
			OwnerId:     course.OwnerId,
		})
	}

	return &courseservicepb.ReadAllCoursesByOwnerIdResponse{
		Courses: coursesPb,
	}, nil
}

func (h *CourseHandler) UpdateCourse(ctx context.Context, req *courseservicepb.UpdateCourseRequest) (*courseservicepb.UpdateCourseResponse, error) {
	updatedCourse, err := h.service.UpdateCourse(ctx, &entities.Course{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		AccessType:  req.AccessType,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &courseservicepb.UpdateCourseResponse{
		Course: &courseservicepb.Course{
			Id:          updatedCourse.Id,
			Title:       updatedCourse.Title,
			Description: updatedCourse.Description,
			AccessType:  updatedCourse.AccessType,
			IsPublished: updatedCourse.IsPublished,
			OwnerId:     updatedCourse.OwnerId,
		},
	}, nil
}

func (h *CourseHandler) DeleteCourse(ctx context.Context, req *courseservicepb.DeleteCourseRequest) (*courseservicepb.DeleteCourseResponse, error) {
	if err := h.service.DeleteCourse(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &courseservicepb.DeleteCourseResponse{}, nil
}
