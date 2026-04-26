package models

import "github.com/google/uuid"

type Furniture struct {
	ID         uuid.UUID
	Name       string
	Code       string
	Price      float64
	StockCount int
}
