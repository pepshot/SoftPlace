package httpserver

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
)

type errorResponse struct {
	Error string `json:"error"`
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, pgx.ErrNoRows):
		return http.StatusNotFound, "not found"
	case errors.Is(err, service.ErrInvalidID):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrInvalidDate):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrInvalidCount):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrEmptyComposition):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrPasswordMismatch):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrLoginAlreadyUsed):
		return http.StatusConflict, err.Error()
	case errors.Is(err, service.ErrEmailAlreadyUsed):
		return http.StatusConflict, err.Error()
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, service.ErrNotEnoughStock):
		return http.StatusConflict, err.Error()
	case errors.Is(err, service.ErrExternalService):
		return http.StatusBadGateway, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
