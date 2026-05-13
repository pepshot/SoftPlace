package models

import "github.com/google/uuid"

type Material struct {
	ID         uuid.UUID
	Name       string
	Code       string
	Price      float64
	StockCount int
}
