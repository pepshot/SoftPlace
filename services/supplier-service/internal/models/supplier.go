package models

import "github.com/google/uuid"

type Supplier struct {
	ID       uuid.UUID
	Login    string
	Password string
	Email    string
}
