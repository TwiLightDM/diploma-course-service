package completed_module_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/internal/errs"
	"github.com/TwiLightDM/diploma-course-service/proto/completedmoduleservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompletedModuleService interface {
	CreateCompletedModule(ctx context.Context, completedModule *entities.CompletedModule) error
	ReadCompletedModuleByUserIdAndModuleId(ctx context.Context, userId string, lessonId string) (*entities.CompletedModule, error)
	ReadAllCompletedModulesByModuleId(ctx context.Context, moduleId string) ([]entities.CompletedModule, error)
	ReadAllCompletedModulesByUserId(ctx context.Context, userId string) ([]entities.CompletedModule, error)
}

type CompletedModuleHandler struct {
	completedmoduleservicepb.UnimplementedCompletedModuleServiceServer
	service CompletedModuleService
}

func NewCompletedModuleHandler(service CompletedModuleService) *CompletedModuleHandler {
	return &CompletedModuleHandler{service: service}
}

func (h *CompletedModuleHandler) CreateCompletedModule(ctx context.Context, req *completedmoduleservicepb.CreateCompletedModuleRequest) (*completedmoduleservicepb.CreateCompletedModuleResponse, error) {
	completedModule := entities.CompletedModule{
		UserId:   req.UserId,
		ModuleId: req.ModuleId,
	}

	err := h.service.CreateCompletedModule(ctx, &completedModule)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateKey) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedmoduleservicepb.CreateCompletedModuleResponse{
		CompletedModule: &completedmoduleservicepb.CompletedModule{
			UserId:   completedModule.UserId,
			ModuleId: completedModule.ModuleId,
		},
	}, nil
}

func (h *CompletedModuleHandler) ReadCompletedModuleByUserIdAndModuleId(ctx context.Context, req *completedmoduleservicepb.ReadCompletedModuleByUserIdAndModuleIdRequest) (*completedmoduleservicepb.ReadCompletedModuleByUserIdAndModuleIdResponse, error) {
	completedModule, err := h.service.ReadCompletedModuleByUserIdAndModuleId(ctx, req.UserId, req.ModuleId)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return &completedmoduleservicepb.ReadCompletedModuleByUserIdAndModuleIdResponse{}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &completedmoduleservicepb.ReadCompletedModuleByUserIdAndModuleIdResponse{
		CompletedModule: &completedmoduleservicepb.CompletedModule{
			UserId:   completedModule.UserId,
			ModuleId: completedModule.ModuleId,
		},
	}, nil
}

func (h *CompletedModuleHandler) ReadAllCompletedModulesByModuleId(ctx context.Context, req *completedmoduleservicepb.ReadAllCompletedModulesByModuleIdRequest) (*completedmoduleservicepb.ReadAllCompletedModulesByModuleIdResponse, error) {
	completedModules, err := h.service.ReadAllCompletedModulesByModuleId(ctx, req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedModulesPb := make([]*completedmoduleservicepb.CompletedModule, 0, len(completedModules))
	for _, completedModule := range completedModules {
		completedModulesPb = append(completedModulesPb, &completedmoduleservicepb.CompletedModule{
			UserId:   completedModule.UserId,
			ModuleId: completedModule.ModuleId,
		},
		)
	}

	return &completedmoduleservicepb.ReadAllCompletedModulesByModuleIdResponse{
		CompletedModules: completedModulesPb,
	}, nil
}

func (h *CompletedModuleHandler) ReadAllCompletedModulesByUserId(ctx context.Context, req *completedmoduleservicepb.ReadAllCompletedModulesByUserIdRequest) (*completedmoduleservicepb.ReadAllCompletedModulesByUserIdResponse, error) {
	completedModules, err := h.service.ReadAllCompletedModulesByUserId(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	completedModulesPb := make([]*completedmoduleservicepb.CompletedModule, 0, len(completedModules))
	for _, completedModule := range completedModules {
		completedModulesPb = append(completedModulesPb, &completedmoduleservicepb.CompletedModule{
			UserId:   completedModule.UserId,
			ModuleId: completedModule.ModuleId,
		},
		)
	}

	return &completedmoduleservicepb.ReadAllCompletedModulesByUserIdResponse{
		CompletedModules: completedModulesPb,
	}, nil
}
