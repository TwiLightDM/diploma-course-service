package lesson_file_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonfileservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LessonFileService interface {
	UploadFile(ctx context.Context, lessonFile *entities.LessonFile, file []byte) error
	GetLessonFiles(ctx context.Context, lessonId string) ([]entities.LessonFile, error)
	DeleteFile(ctx context.Context, id, objectName string) error
}

type LessonFileHandler struct {
	lessonfileservicepb.UnimplementedLessonFileServiceServer
	service LessonFileService
}

func NewLessonFileHandler(service LessonFileService) *LessonFileHandler {
	return &LessonFileHandler{service: service}
}

func (h *LessonFileHandler) UploadFile(ctx context.Context, req *lessonfileservicepb.UploadFileRequest) (*lessonfileservicepb.UploadFileResponse, error) {
	lessonFile := entities.LessonFile{
		FileName:    req.FileName,
		ContentType: req.ContentType,
		LessonId:    req.LessonId,
		Size:        req.Size,
	}
	err := h.service.UploadFile(ctx, &lessonFile, req.File)
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

func (h *LessonFileHandler) GetLessonFiles(ctx context.Context, req *lessonfileservicepb.GetLessonFilesRequest) (*lessonfileservicepb.GetLessonFilesResponse, error) {
	lessonFiles, err := h.service.GetLessonFiles(ctx, req.LessonId)
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

func (h *LessonFileHandler) DeleteFile(ctx context.Context, req *lessonfileservicepb.DeleteFileRequest) (*lessonfileservicepb.DeleteFileResponse, error) {
	if err := h.service.DeleteFile(ctx, req.Id, req.ObjectName); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lessonfileservicepb.DeleteFileResponse{}, nil
}
