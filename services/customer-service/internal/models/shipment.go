package models

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	Id         uuid.UUID
	CustomerID uuid.UUID
	Code       string
	Date       time.Time
	Price      float64
}
