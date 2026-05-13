package model

import "github.com/google/uuid"

type Garniture struct {
	ID         uuid.UUID
	Name       string
	Code       string
	Price      float64
	StockCount int
}
