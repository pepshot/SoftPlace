package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	supplierpb "github.com/pepshot/SoftPlace/gen/supplier"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/shared/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	supplierpb.UnimplementedSupplierServiceServer
	moduleService service.ModuleUseCase
	logger        *logger.Logger
}

func NewServer(moduleService service.ModuleUseCase, logger *logger.Logger) *Server {
	return &Server{
		moduleService: moduleService,
		logger:        logger,
	}
}

func (s *Server) GetModule(ctx context.Context, req *supplierpb.GetModuleRequest) (*supplierpb.ModuleResponse, error) {
	if req == nil || req.GetModuleId() == "" {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidID.Error())
	}

	moduleID, err := uuid.Parse(req.GetModuleId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidID.Error())
	}

	module, err := s.moduleService.GetByID(ctx, moduleID)
	if err != nil {
		return nil, toStatusError(err)
	}

	return moduleResponseFromDTO(module), nil
}

func (s *Server) ReserveModules(ctx context.Context, req *supplierpb.ReserveModulesRequest) (*supplierpb.ReserveModulesResponse, error) {
	items := make([]service.ModuleItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item == nil {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidComposition.Error())
		}

		if item.GetModuleId() == "" {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidID.Error())
		}

		if item.GetCount() <= 0 {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidCount.Error())
		}

		items = append(items, service.ModuleItem{
			ModuleID: item.GetModuleId(),
			Count:    int(item.GetCount()),
		})
	}

	if err := s.moduleService.ReserveModules(ctx, items); err != nil {
		return nil, toStatusError(err)
	}

	return &supplierpb.ReserveModulesResponse{Success: true, Message: "modules reserved successfully"}, nil
}

func (s *Server) ReleaseModules(ctx context.Context, req *supplierpb.ReleaseModulesRequest) (*supplierpb.ReleaseModulesResponse, error) {
	items := make([]service.ModuleItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item == nil {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidComposition.Error())
		}

		if item.GetModuleId() == "" {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidID.Error())
		}

		if item.GetCount() <= 0 {
			return nil, status.Error(codes.InvalidArgument, service.ErrInvalidCount.Error())
		}

		items = append(items, service.ModuleItem{
			ModuleID: item.GetModuleId(),
			Count:    int(item.GetCount()),
		})
	}

	if err := s.moduleService.ReleaseModules(ctx, items); err != nil {
		return nil, toStatusError(err)
	}

	return &supplierpb.ReleaseModulesResponse{Success: true, Message: "modules released successfully"}, nil
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrInvalidID), errors.Is(err, service.ErrInvalidCount), errors.Is(err, service.ErrInvalidPrice), errors.Is(err, service.ErrInvalidComposition):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotEnoughStock):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, service.ErrCodeAlreadyUsed):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func moduleResponseFromDTO(item dto.ModuleResponse) *supplierpb.ModuleResponse {
	return &supplierpb.ModuleResponse{
		Id:         item.ID,
		Name:       item.Name,
		Code:       item.Code,
		Price:      item.Price,
		StockCount: int32(item.StockCount),
	}
}
