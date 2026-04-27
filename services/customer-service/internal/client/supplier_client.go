package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	supplierpb "github.com/pepshot/SoftPlace/gen/supplier"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
	"github.com/pepshot/SoftPlace/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultRequestTimeout = 5 * time.Second

type SupplierGrpcClient struct {
	conn   *grpc.ClientConn
	client supplierpb.SupplierServiceClient
	logger *logger.Logger
}

func NewSupplierGrpcClient(address string, logger *logger.Logger) (*SupplierGrpcClient, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to supplier-service: %w", err)
	}

	return &SupplierGrpcClient{
		conn:   conn,
		client: supplierpb.NewSupplierServiceClient(conn),
		logger: logger,
	}, nil
}

func (c *SupplierGrpcClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *SupplierGrpcClient) GetModule(ctx context.Context, moduleID string) (service.ModuleInfo, error) {
	c.logger.Debug("sending grpc GetModule request", "moduleID", moduleID)

	if moduleID == "" {
		return service.ModuleInfo{}, errors.New("module id is empty")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	response, err := c.client.GetModule(ctx, &supplierpb.GetModuleRequest{
		ModuleId: moduleID,
	})
	if err != nil {
		c.logger.Error(
			"grpc GetModule request failed",
			"moduleID", moduleID,
			"error", err,
		)

		return service.ModuleInfo{}, err
	}

	module := service.ModuleInfo{
		ID:         response.Id,
		Name:       response.Name,
		Code:       response.Code,
		Price:      response.Price,
		StockCount: int(response.StockCount),
	}

	c.logger.Debug(
		"grpc GetModule request completed",
		"moduleID", module.ID,
		"code", module.Code,
		"stockCount", module.StockCount,
	)

	return module, nil
}

func (c *SupplierGrpcClient) ReserveModules(ctx context.Context, items []service.ModuleItem) error {
	c.logger.Debug("sending grpc ReserveModules request", "itemsCount", len(items))

	if len(items) == 0 {
		return errors.New("module items cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	request := &supplierpb.ReserveModulesRequest{
		Items: make([]*supplierpb.ModuleItem, 0, len(items)),
	}

	for _, item := range items {
		if item.ModuleID == "" {
			return errors.New("module id is empty")
		}

		if item.Count <= 0 {
			return errors.New("module count must be greater than zero")
		}

		request.Items = append(request.Items, &supplierpb.ModuleItem{
			ModuleId: item.ModuleID,
			Count:    int32(item.Count),
		})
	}

	response, err := c.client.ReserveModules(ctx, request)
	if err != nil {
		c.logger.Error(
			"grpc ReserveModules request failed",
			"itemsCount", len(items),
			"error", err,
		)

		return err
	}

	if !response.Success {
		c.logger.Warn(
			"grpc ReserveModules request rejected",
			"message", response.Message,
		)

		if response.Message == "" {
			return errors.New("reserve modules request rejected")
		}

		return errors.New(response.Message)
	}

	c.logger.Debug("grpc ReserveModules request completed", "itemsCount", len(items))

	return nil
}

func (c *SupplierGrpcClient) ReleaseModules(ctx context.Context, items []service.ModuleItem) error {
	c.logger.Debug("sending grpc ReleaseModules request", "itemsCount", len(items))

	if len(items) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	request := &supplierpb.ReleaseModulesRequest{
		Items: make([]*supplierpb.ModuleItem, 0, len(items)),
	}

	for _, item := range items {
		if item.ModuleID == "" {
			return errors.New("module id is empty")
		}

		if item.Count <= 0 {
			return errors.New("module count must be greater than zero")
		}

		request.Items = append(request.Items, &supplierpb.ModuleItem{
			ModuleId: item.ModuleID,
			Count:    int32(item.Count),
		})
	}

	response, err := c.client.ReleaseModules(ctx, request)
	if err != nil {
		c.logger.Error(
			"grpc ReleaseModules request failed",
			"itemsCount", len(items),
			"error", err,
		)

		return err
	}

	if !response.Success {
		c.logger.Warn(
			"grpc ReleaseModules request rejected",
			"message", response.Message,
		)

		if response.Message == "" {
			return errors.New("release modules request rejected")
		}

		return errors.New(response.Message)
	}

	c.logger.Debug("grpc ReleaseModules request completed", "itemsCount", len(items))

	return nil
}
