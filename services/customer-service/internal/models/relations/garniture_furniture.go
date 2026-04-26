package relations

import "github.com/google/uuid"

type GarnitureFurniture struct {
	GarnitureID uuid.UUID
	FurnitureID uuid.UUID
	Count       int
}
