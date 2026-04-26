package relations

import "github.com/google/uuid"

type ModuleMaterial struct {
	ModuleID   uuid.UUID
	MaterialID uuid.UUID
	Count      int
}
