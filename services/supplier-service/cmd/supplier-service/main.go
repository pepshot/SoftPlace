package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	supplierpb "github.com/pepshot/SoftPlace/gen/supplier"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	grpcserver "github.com/pepshot/SoftPlace/services/supplier-service/internal/transport/grpc"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/database"
	"github.com/pepshot/SoftPlace/shared/logger"
	"google.golang.org/grpc"
)

func main() {
	envPath := os.Getenv("CONFIG_PATH")
	if strings.TrimSpace(envPath) == "" {
		envPath = filepath.Join("configs", "supplier.env")
	}

	cfg := config.LoadConfig(envPath)

	log, err := logger.New(cfg.Logger, cfg.App.ServiceName)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, cfg.Database, log)
	if err != nil {
		log.Error("failed to init postgres", "error", err)
		return
	}
	defer db.Close()

	moduleRepo := repository.NewModuleRepository(db, log)
	moduleMaterialRepo := repository.NewModuleMaterialRepository(db, log)
	materialRepo := repository.NewMaterialRepository(db, log)
	txManager := repository.NewTxManager(db, log)
	moduleService := service.NewModuleService(moduleRepo, moduleMaterialRepo, materialRepo, txManager, log)

	grpcServer := grpc.NewServer()
	supplierpb.RegisterSupplierServiceServer(grpcServer, grpcserver.NewServer(moduleService, log))

	listener, err := net.Listen("tcp", cfg.GRPC.Host+":"+cfg.GRPC.Port)
	if err != nil {
		log.Error("failed to listen grpc", "error", err)
		return
	}

	go func() {
		log.Info("supplier-service gRPC started", "addr", listener.Addr().String())

		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error("grpc server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		grpcServer.Stop()
	}
}
