package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/services/supplier-service/tests/mocks"
)

func TestMaterialServiceGetListSuccess(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))

	repo.On("GetList", mock.Anything).Return([]models.Material{sampleMaterial(uuid.New())}, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestMaterialServiceGetByIDNotFound(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(models.Material{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), id)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestMaterialServiceCreateInvalidPrice(t *testing.T) {
	svc := service.NewMaterialService(&mocks.MaterialRepositoryMock{}, testLogger(t))
	req := materialRequest()
	req.Price = -1

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidPrice)
}

func TestMaterialServiceCreateCodeAlreadyUsed(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))
	req := materialRequest()
	repo.On("ExistsByCode", mock.Anything, req.Code).Return(true, nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrCodeAlreadyUsed)
}

func TestMaterialServiceCreateSuccess(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))
	req := materialRequest()

	repo.On("ExistsByCode", mock.Anything, req.Code).Return(false, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("models.Material")).Return(nil)

	id, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestMaterialServiceUpdateNotFound(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))
	id := uuid.New()
	req := materialRequest()

	repo.On("GetByID", mock.Anything, id).Return(models.Material{}, errors.New("not found"))

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestMaterialServiceDeleteNotFound(t *testing.T) {
	repo := &mocks.MaterialRepositoryMock{}
	svc := service.NewMaterialService(repo, testLogger(t))
	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(errors.New("not found"))

	err := svc.Delete(context.Background(), id)
	require.ErrorIs(t, err, service.ErrNotFound)
}
