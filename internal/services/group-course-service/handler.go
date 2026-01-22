package group_course_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/groupcourseservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupCourseHandler struct {
	groupcourseservicepb.UnimplementedGroupCourseServiceServer
	service GroupCourseService
}

func NewGroupCourseHandler(service GroupCourseService) *GroupCourseHandler {
	return &GroupCourseHandler{service: service}
}

func (h *GroupCourseHandler) CreateGroupCourse(ctx context.Context, req *groupcourseservicepb.CreateGroupCourseRequest) (*groupcourseservicepb.CreateGroupCourseResponse, error) {
	groupCourse, err := h.service.CreateGroupCourse(ctx, &entities.GroupCourse{
		CourseId: req.CourseId,
		GroupId:  req.GroupId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupcourseservicepb.CreateGroupCourseResponse{
		GroupCourse: &groupcourseservicepb.GroupCourse{
			Id:       groupCourse.Id,
			CourseId: groupCourse.CourseId,
			GroupId:  groupCourse.GroupId,
		},
	}, nil
}

func (h *GroupCourseHandler) ReadAllGroupCoursesByCourseId(ctx context.Context, req *groupcourseservicepb.ReadAllGroupCoursesByCourseIdRequest) (*groupcourseservicepb.ReadAllGroupCoursesByCourseIdResponse, error) {
	groupCourses, err := h.service.ReadAllGroupCoursesByCourseId(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	groupCoursesPb := h.groupCoursesToPb(groupCourses)

	return &groupcourseservicepb.ReadAllGroupCoursesByCourseIdResponse{
		GroupCourses: groupCoursesPb,
	}, nil
}

func (h *GroupCourseHandler) ReadAllGroupCoursesByGroupId(ctx context.Context, req *groupcourseservicepb.ReadAllGroupCoursesByGroupIdRequest) (*groupcourseservicepb.ReadAllGroupCoursesByGroupIdResponse, error) {
	groupCourses, err := h.service.ReadAllGroupCoursesByGroupId(ctx, req.GroupId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	groupCoursesPb := h.groupCoursesToPb(groupCourses)

	return &groupcourseservicepb.ReadAllGroupCoursesByGroupIdResponse{
		GroupCourses: groupCoursesPb,
	}, nil
}

func (h *GroupCourseHandler) DeleteGroupCourse(ctx context.Context, req *groupcourseservicepb.DeleteGroupCourseRequest) (*groupcourseservicepb.DeleteGroupCourseResponse, error) {
	err := h.service.DeleteGroupCourse(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupcourseservicepb.DeleteGroupCourseResponse{}, nil
}

func (h *GroupCourseHandler) groupCoursesToPb(groupCourses []entities.GroupCourse) []*groupcourseservicepb.GroupCourse {
	groupCoursesPb := make([]*groupcourseservicepb.GroupCourse, 0, len(groupCourses))
	for _, groupCourse := range groupCourses {
		groupCoursesPb = append(groupCoursesPb, &groupcourseservicepb.GroupCourse{
			Id:       groupCourse.Id,
			CourseId: groupCourse.CourseId,
			GroupId:  groupCourse.GroupId,
		})
	}

	return groupCoursesPb
}
