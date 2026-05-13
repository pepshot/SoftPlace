package model

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Code       string
	Date       time.Time
	Price      float64
}
