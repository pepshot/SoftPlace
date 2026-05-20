package httpserver

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
)

type errorResponse struct {
	Error string `json:"error"`
}

// mapServiceError maps domain errors to HTTP status codes
func mapServiceError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, service.ErrNotFound), errors.Is(err, pgx.ErrNoRows):
		return http.StatusNotFound
	case errors.Is(err, service.ErrInvalidID), errors.Is(err, service.ErrInvalidCount), errors.Is(err, service.ErrInvalidPrice),
		errors.Is(err, service.ErrInvalidComposition), errors.Is(err, service.ErrInvalidDate), errors.Is(err, service.ErrEmptyComposition),
		errors.Is(err, service.ErrPasswordMismatch):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrNotEnoughStock), errors.Is(err, service.ErrCodeAlreadyUsed), errors.Is(err, service.ErrLoginAlreadyUsed), errors.Is(err, service.ErrEmailAlreadyUsed):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
