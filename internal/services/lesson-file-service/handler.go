package lesson_file_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonfileservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LessonFileService interface {
	CreateLessonFile(ctx context.Context, lessonFile *entities.LessonFile, file []byte) error
	ReadAllLessonFilesByLessonId(ctx context.Context, lessonId string) ([]entities.LessonFile, error)
	DeleteLessonFile(ctx context.Context, id string) error
}

type LessonFileHandler struct {
	lessonfileservicepb.UnimplementedLessonFileServiceServer
	service LessonFileService
}

func NewLessonFileHandler(service LessonFileService) *LessonFileHandler {
	return &LessonFileHandler{service: service}
}

func (h *LessonFileHandler) CreateLessonFile(ctx context.Context, req *lessonfileservicepb.UploadFileRequest) (*lessonfileservicepb.UploadFileResponse, error) {
	lessonFile := entities.LessonFile{
		FileName:    req.FileName,
		ContentType: req.ContentType,
		LessonId:    req.LessonId,
	}
	err := h.service.CreateLessonFile(ctx, &lessonFile, req.File)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonfileservicepb.UploadFileResponse{
		File: &lessonfileservicepb.LessonFile{
			Id:         lessonFile.Id,
			ObjectName: lessonFile.ObjectName,
			Url:        lessonFile.Url,
		},
	}, nil
}

func (h *LessonFileHandler) ReadAllLessonFilesByModuleId(ctx context.Context, req *lessonfileservicepb.GetLessonFilesRequest) (*lessonfileservicepb.GetLessonFilesResponse, error) {
	lessonFiles, err := h.service.ReadAllLessonFilesByLessonId(ctx, req.LessonId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	lessonFilesPb := make([]*lessonfileservicepb.LessonFile, 0, len(lessonFiles))
	for _, lessonFile := range lessonFiles {
		lessonFilesPb = append(lessonFilesPb, &lessonfileservicepb.LessonFile{
			Id:         lessonFile.Id,
			ObjectName: lessonFile.ObjectName,
			Url:        lessonFile.Url,
		})
	}

	return &lessonfileservicepb.GetLessonFilesResponse{
		Files: lessonFilesPb,
	}, nil
}

func (h *LessonFileHandler) DeleteLessonFile(ctx context.Context, req *lessonfileservicepb.DeleteFileRequest) (*lessonfileservicepb.DeleteFileResponse, error) {
	if err := h.service.DeleteLessonFile(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonfileservicepb.DeleteFileResponse{}, nil
}
