package relations

import "github.com/google/uuid"

type SupplyMaterial struct {
	SupplyID   uuid.UUID
	MaterialID uuid.UUID
	Count      int
}
