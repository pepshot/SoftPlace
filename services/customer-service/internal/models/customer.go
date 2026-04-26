package models

import "github.com/google/uuid"

type Customer struct {
	ID       uuid.UUID
	Login    string
	Email    string
	Password string
}
