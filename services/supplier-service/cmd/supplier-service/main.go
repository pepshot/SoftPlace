package main

// @title Supplier Service API
// @version 1.0
// @description Supplier service API
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"errors"
	"net"
	"net/http"
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
	httpserver "github.com/pepshot/SoftPlace/services/supplier-service/internal/transport/http"
	"github.com/pepshot/SoftPlace/shared/auth"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/database"
	"github.com/pepshot/SoftPlace/shared/logger"
	"google.golang.org/grpc"

	_ "github.com/pepshot/SoftPlace/services/supplier-service/docs"
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
	supplyRepo := repository.NewSupplyRepository(db, log)
	supplyMaterialRepo := repository.NewSupplyMaterialRepository(db, log)
	supplierRepo := repository.NewSupplierRepository(db)
	txManager := repository.NewTxManager(db, log)

	jwtManager := auth.NewJWTManager(cfg.Auth)

	moduleService := service.NewModuleService(moduleRepo, moduleMaterialRepo, materialRepo, txManager, log)
	materialService := service.NewMaterialService(materialRepo, log)
	supplyService := service.NewSupplyService(supplyRepo, supplyMaterialRepo, materialRepo, txManager, log)
	supplierService := service.NewSupplierService(supplierRepo, jwtManager, log)

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

	// HTTP server
	handler := httpserver.NewHandler(moduleService, materialService, supplyService, supplierService, log)
	router := httpserver.NewRouter(handler, jwtManager)

	httpServer := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler: router,
	}

	go func() {
		if err := startServer(httpServer, cfg.HTTP, log); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "error", err)
		}
	}()

	log.Info("supplier-service HTTP started", "addr", httpServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown failed", "error", err)
	}

	// Shutdown gRPC server
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

func startServer(server *http.Server, cfg config.HTTPConfig, log *logger.Logger) error {
	if fileExists(cfg.CertFile) && fileExists(cfg.KeyFile) {
		log.Info("starting HTTPS server", "addr", server.Addr)
		return server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	}

	log.Info("starting HTTP server", "addr", server.Addr)
	return server.ListenAndServe()
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
