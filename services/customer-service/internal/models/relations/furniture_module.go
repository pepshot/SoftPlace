package relations

import "github.com/google/uuid"

type FurnitureModule struct {
	FurnitureID uuid.UUID
	ModuleId    uuid.UUID
	Count       int
}
