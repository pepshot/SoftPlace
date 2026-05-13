package main

// @title Customer Service API
// @version 1.0
// @description Customer service API
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/client"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
	httpserver "github.com/pepshot/SoftPlace/services/customer-service/internal/transport/http"
	"github.com/pepshot/SoftPlace/shared/auth"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/database"
	"github.com/pepshot/SoftPlace/shared/logger"

	_ "github.com/pepshot/SoftPlace/services/customer-service/docs"
)

func main() {
	envPath := os.Getenv("CONFIG_PATH")
	if strings.TrimSpace(envPath) == "" {
		envPath = filepath.Join("configs", "customer.env")
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

	supplierClient, err := client.NewSupplierGrpcClient(cfg.GRPC.SupplierServiceAddr, log)
	if err != nil {
		log.Error("failed to init supplier client", "error", err)
		return
	}
	defer func() { _ = supplierClient.Close() }()

	customerRepo := repository.NewCustomerRepository(db, log)
	furnitureRepo := repository.NewFurnitureRepository(db, log)
	garnitureRepo := repository.NewGarnitureRepository(db, log)
	shipmentRepo := repository.NewShipmentRepository(db, log)
	furnitureModuleRepo := repository.NewFurnitureModuleRepository(db, log)
	garnitureFurnitureRepo := repository.NewGarnitureFurnitureRepository(db, log)
	shipmentGarnitureRepo := repository.NewShipmentGarnitureRepository(db, log)

	txManager := repository.NewTxManager(db, log)
	jwtManager := auth.NewJWTManager(cfg.Auth)

	customerService := service.NewCustomerService(customerRepo, jwtManager, log)
	furnitureService := service.NewFurnitureService(&furnitureRepo, furnitureModuleRepo, supplierClient, txManager, log)
	garnitureService := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, &furnitureRepo, txManager, log)
	shipmentService := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, log)

	handler := httpserver.NewHandler(customerService, furnitureService, garnitureService, shipmentService, log)
	router := httpserver.NewRouter(handler, jwtManager)

	server := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler: router,
	}

	go func() {
		if err := startServer(server, cfg.HTTP, log); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
		}
	}()

	log.Info("customer-service started", "addr", server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", "error", err)
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
