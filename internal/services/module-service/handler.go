package module_service

import (
	"context"
	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/TwiLightDM/diploma-course-service/proto/moduleservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModuleHandler struct {
	moduleservicepb.UnimplementedModuleServiceServer
	service ModuleService
}

func NewModuleHandler(service ModuleService) *ModuleHandler {
	return &ModuleHandler{service: service}
}

func (h *ModuleHandler) CreateModule(ctx context.Context, req *moduleservicepb.CreateModuleRequest) (*moduleservicepb.CreateModuleResponse, error) {
	module := entities.Module{
		Title:       req.Title,
		Description: req.Description,
		CourseId:    req.CourseId,
	}

	err := h.service.CreateModule(ctx, &module)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moduleservicepb.CreateModuleResponse{
		Module: &moduleservicepb.Module{
			Id:          module.Id,
			Title:       module.Title,
			Description: module.Description,
			CourseId:    module.CourseId,
		},
	}, nil
}

func (h *ModuleHandler) ReadModule(ctx context.Context, req *moduleservicepb.ReadModuleRequest) (*moduleservicepb.ReadModuleResponse, error) {
	module, err := h.service.ReadModuleById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moduleservicepb.ReadModuleResponse{
		Module: &moduleservicepb.Module{
			Id:          module.Id,
			Title:       module.Title,
			Description: module.Description,
			Position:    module.Position,
			CourseId:    module.CourseId,
		},
	}, nil
}

func (h *ModuleHandler) ReadAllModulesByCourseId(ctx context.Context, req *moduleservicepb.ReadAllModulesByCourseIdRequest) (*moduleservicepb.ReadAllModulesByCourseIdResponse, error) {
	modules, err := h.service.ReadAllModulesByCourseId(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	modulesPb := make([]*moduleservicepb.Module, 0, len(modules))
	for _, module := range modules {
		modulesPb = append(modulesPb, &moduleservicepb.Module{
			Id:          module.Id,
			Title:       module.Title,
			Description: module.Description,
			Position:    module.Position,
			CourseId:    module.CourseId,
		})
	}

	return &moduleservicepb.ReadAllModulesByCourseIdResponse{
		Modules: modulesPb,
	}, nil
}

func (h *ModuleHandler) UpdateModule(ctx context.Context, req *moduleservicepb.UpdateModuleRequest) (*moduleservicepb.UpdateModuleResponse, error) {
	updatedModule, err := h.service.UpdateModule(ctx, &entities.Module{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moduleservicepb.UpdateModuleResponse{
		Module: &moduleservicepb.Module{
			Id:          updatedModule.Id,
			Title:       updatedModule.Title,
			Description: updatedModule.Description,
			Position:    updatedModule.Position,
			CourseId:    updatedModule.CourseId,
		},
	}, nil
}

func (h *ModuleHandler) DeleteModule(ctx context.Context, req *moduleservicepb.DeleteModuleRequest) (*moduleservicepb.DeleteModuleResponse, error) {
	if err := h.service.DeleteModule(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &moduleservicepb.DeleteModuleResponse{}, nil
}
