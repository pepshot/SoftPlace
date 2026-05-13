package models

import "github.com/google/uuid"

type Module struct {
	ID         uuid.UUID
	Name       string
	Code       string
	Price      float64
	StockCount int
}
