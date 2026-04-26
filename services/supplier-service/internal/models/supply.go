package models

import (
	"time"

	"github.com/google/uuid"
)

type Supply struct {
	ID         uuid.UUID
	SupplierID uuid.UUID
	Code       string
	Date       time.Time
	Price      float64
}
