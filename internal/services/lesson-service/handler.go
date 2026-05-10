package lesson_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LessonService interface {
	CreateLesson(ctx context.Context, lesson *entities.Lesson) error
	ReadLessonById(ctx context.Context, id string) (*entities.Lesson, error)
	ReadAllLessonsByModuleId(ctx context.Context, moduleId string) ([]entities.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *entities.Lesson) (*entities.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
}

type LessonHandler struct {
	lessonservicepb.UnimplementedLessonServiceServer
	service LessonService
}

func NewLessonHandler(service LessonService) *LessonHandler {
	return &LessonHandler{service: service}
}

func (h *LessonHandler) CreateLesson(ctx context.Context, req *lessonservicepb.CreateLessonRequest) (*lessonservicepb.CreateLessonResponse, error) {
	lesson := entities.Lesson{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		ModuleId:    req.ModuleId,
	}

	err := h.service.CreateLesson(ctx, &lesson)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonservicepb.CreateLessonResponse{
		Lesson: &lessonservicepb.Lesson{
			Id:          lesson.Id,
			Title:       lesson.Title,
			Description: lesson.Description,
			Position:    lesson.Position,
			Content:     lesson.Content,
			ModuleId:    lesson.ModuleId,
		},
	}, nil
}

func (h *LessonHandler) ReadLesson(ctx context.Context, req *lessonservicepb.ReadLessonRequest) (*lessonservicepb.ReadLessonResponse, error) {
	lesson, err := h.service.ReadLessonById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	files := make([]*lessonservicepb.LessonFile, 0, len(lesson.Files))
	for _, f := range lesson.Files {
		files = append(files, &lessonservicepb.LessonFile{
			Id:         f.Id,
			ObjectName: f.ObjectName,
			Url:        f.Url,
		})
	}

	return &lessonservicepb.ReadLessonResponse{
		Lesson: &lessonservicepb.Lesson{
			Id:          lesson.Id,
			Title:       lesson.Title,
			Description: lesson.Description,
			Position:    lesson.Position,
			Content:     lesson.Content,
			ModuleId:    lesson.ModuleId,
			Files:       files,
		},
	}, nil
}

func (h *LessonHandler) ReadAllLessonsByModuleId(ctx context.Context, req *lessonservicepb.ReadAllLessonsByModuleIdRequest) (*lessonservicepb.ReadAllLessonsByModuleIdResponse, error) {
	lessons, err := h.service.ReadAllLessonsByModuleId(ctx, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	lessonsPb := make([]*lessonservicepb.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		files := make([]*lessonservicepb.LessonFile, 0, len(lesson.Files))

		for _, f := range lesson.Files {
			files = append(files, &lessonservicepb.LessonFile{
				Id:         f.Id,
				ObjectName: f.ObjectName,
				Url:        f.Url,
			})
		}

		lessonsPb = append(lessonsPb, &lessonservicepb.Lesson{
			Id:          lesson.Id,
			Title:       lesson.Title,
			Description: lesson.Description,
			Position:    lesson.Position,
			Content:     lesson.Content,
			ModuleId:    lesson.ModuleId,
			Files:       files,
		})
	}

	return &lessonservicepb.ReadAllLessonsByModuleIdResponse{
		Lessons: lessonsPb,
	}, nil
}

func (h *LessonHandler) UpdateLesson(ctx context.Context, req *lessonservicepb.UpdateLessonRequest) (*lessonservicepb.UpdateLessonResponse, error) {
	updatedLesson, err := h.service.UpdateLesson(ctx, &entities.Lesson{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonservicepb.UpdateLessonResponse{
		Lesson: &lessonservicepb.Lesson{
			Id:          updatedLesson.Id,
			Title:       updatedLesson.Title,
			Description: updatedLesson.Description,
			Position:    updatedLesson.Position,
			Content:     updatedLesson.Content,
			ModuleId:    updatedLesson.ModuleId,
		},
	}, nil
}

func (h *LessonHandler) DeleteLesson(ctx context.Context, req *lessonservicepb.DeleteLessonRequest) (*lessonservicepb.DeleteLessonResponse, error) {
	if err := h.service.DeleteLesson(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonservicepb.DeleteLessonResponse{}, nil
}
