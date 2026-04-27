package service

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type CustomerService struct {
	customerRepo CustomerRepository
	jwtUseCase   JWTUseCase
	logger       *logger.Logger
}

func NewCustomerService(customerRepo CustomerRepository, jwtUseCase JWTUseCase, logger *logger.Logger) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		jwtUseCase:   jwtUseCase,
		logger:       logger,
	}
}

func (s *CustomerService) Register(ctx context.Context, req dto.RegisterCustomerRequest) (dto.AuthCustomerResponse, error) {
	s.logger.Info("starting customer registration", "login", req.Login, "email", req.Email)

	if req.Password != req.ConfirmPassword {
		s.logger.Warn("customer registration failed: password mismatch", "login", req.Login)
		return dto.AuthCustomerResponse{}, ErrPasswordMismatch
	}

	existsLogin, err := s.customerRepo.ExistsByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Error("failed to check customer login", "login", req.Login, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	if existsLogin {
		s.logger.Warn("customer registration failed: login already used", "login", req.Login)
		return dto.AuthCustomerResponse{}, ErrLoginAlreadyUsed
	}

	existsEmail, err := s.customerRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("failed to check customer email", "email", req.Email, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	if existsEmail {
		s.logger.Warn("customer registration failed: email already used", "email", req.Email)
		return dto.AuthCustomerResponse{}, ErrEmailAlreadyUsed
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash customer password", "login", req.Login, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	customer := model.Customer{
		ID:       uuid.New(),
		Login:    req.Login,
		Password: string(passwordHash),
		Email:    req.Email,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		s.logger.Error("failed to create customer", "login", req.Login, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	token, err := s.jwtUseCase.GenerateToken(customer.ID.String(), "customer")
	if err != nil {
		s.logger.Error("failed to generate customer token", "customerID", customer.ID, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	response := dto.AuthCustomerResponse{
		Token: token,
		Customer: dto.CustomerResponse{
			ID:    customer.ID.String(),
			Login: customer.Login,
			Email: customer.Email,
		},
	}

	s.logger.Info("customer registered successfully", "customerID", customer.ID)

	return response, nil
}

func (s *CustomerService) Login(ctx context.Context, req dto.LoginCustomerRequest) (dto.AuthCustomerResponse, error) {
	s.logger.Info("starting customer login", "login", req.Login)

	customer, err := s.customerRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Warn("customer login failed: customer not found", "login", req.Login)
		return dto.AuthCustomerResponse{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(customer.Password),
		[]byte(req.Password),
	)
	if err != nil {
		s.logger.Warn("customer login failed: invalid password", "login", req.Login)
		return dto.AuthCustomerResponse{}, ErrInvalidCredentials
	}

	token, err := s.jwtUseCase.GenerateToken(customer.ID.String(), "customer")
	if err != nil {
		s.logger.Error("failed to generate customer token", "customerID", customer.ID, "error", err)
		return dto.AuthCustomerResponse{}, err
	}

	response := dto.AuthCustomerResponse{
		Token:    token,
		Customer: mapper.ToCustomerResponse(customer),
	}

	s.logger.Info("customer logged in successfully", "customerID", customer.ID)

	return response, nil
}

func (s *CustomerService) GetProfile(ctx context.Context, customerID uuid.UUID) (dto.CustomerResponse, error) {
	s.logger.Debug("getting customer profile", "customerID", customerID)

	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		s.logger.Error("failed to get customer profile", "customerID", customerID, "error", err)
		return dto.CustomerResponse{}, ErrNotFound
	}

	return mapper.ToCustomerResponse(customer), nil
}
