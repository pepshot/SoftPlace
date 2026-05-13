package relations

import "github.com/google/uuid"

type ShipmentGarniture struct {
	ShipmentID  uuid.UUID
	GarnitureID uuid.UUID
	Count       int
}
