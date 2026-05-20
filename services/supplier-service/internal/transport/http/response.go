package httpserver

import (
	"net/http"

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
	case err == service.ErrNotFound:
		return http.StatusNotFound
	case err == service.ErrInvalidID, err == service.ErrInvalidCount, err == service.ErrInvalidPrice, err == service.ErrInvalidComposition, err == service.ErrInvalidDate, err == service.ErrEmptyComposition, err == service.ErrPasswordMismatch, err == service.ErrInvalidCredentials:
		return http.StatusBadRequest
	case err == service.ErrNotEnoughStock:
		return http.StatusPreconditionFailed
	case err == service.ErrCodeAlreadyUsed, err == service.ErrLoginAlreadyUsed, err == service.ErrEmailAlreadyUsed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
