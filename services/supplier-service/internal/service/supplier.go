package service

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type SupplierService struct {
	supplierRepo SupplierRepository
	jwtUseCase   JWTUseCase
	logger       *logger.Logger
}

func NewSupplierService(supplierRepo SupplierRepository, jwtUseCase JWTUseCase, logger *logger.Logger) *SupplierService {
	return &SupplierService{supplierRepo: supplierRepo, jwtUseCase: jwtUseCase, logger: logger}
}

func (s *SupplierService) Register(ctx context.Context, req dto.RegisterSupplierRequest) (dto.AuthSupplierResponse, error) {
	if req.Password != req.ConfirmPassword {
		return dto.AuthSupplierResponse{}, ErrPasswordMismatch
	}

	existsLogin, err := s.supplierRepo.ExistsByLogin(ctx, req.Login)
	if err != nil {
		return dto.AuthSupplierResponse{}, err
	}
	if existsLogin {
		return dto.AuthSupplierResponse{}, ErrLoginAlreadyUsed
	}

	existsEmail, err := s.supplierRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthSupplierResponse{}, err
	}
	if existsEmail {
		return dto.AuthSupplierResponse{}, ErrEmailAlreadyUsed
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthSupplierResponse{}, err
	}

	supplier := models.Supplier{
		ID:       uuid.New(),
		Login:    req.Login,
		Password: string(passwordHash),
		Email:    req.Email,
	}

	if err := s.supplierRepo.Create(ctx, supplier); err != nil {
		return dto.AuthSupplierResponse{}, err
	}

	token, err := s.jwtUseCase.GenerateToken(supplier.ID.String(), "supplier")
	if err != nil {
		return dto.AuthSupplierResponse{}, err
	}

	return dto.AuthSupplierResponse{Token: token, Supplier: mapper.ToSupplierResponse(supplier)}, nil
}

func (s *SupplierService) Login(ctx context.Context, req dto.LoginSupplierRequest) (dto.AuthSupplierResponse, error) {
	supplier, err := s.supplierRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		return dto.AuthSupplierResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(supplier.Password), []byte(req.Password)); err != nil {
		return dto.AuthSupplierResponse{}, ErrInvalidCredentials
	}

	token, err := s.jwtUseCase.GenerateToken(supplier.ID.String(), "supplier")
	if err != nil {
		return dto.AuthSupplierResponse{}, err
	}

	return dto.AuthSupplierResponse{Token: token, Supplier: mapper.ToSupplierResponse(supplier)}, nil
}

func (s *SupplierService) GetProfile(ctx context.Context, supplierID uuid.UUID) (dto.SupplierResponse, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, supplierID)
	if err != nil {
		return dto.SupplierResponse{}, ErrNotFound
	}

	return mapper.ToSupplierResponse(supplier), nil
}
