package relations

import "github.com/google/uuid"

type FurnitureModule struct {
	FurnitureID uuid.UUID
	ModuleID    uuid.UUID
	Count       int
}
