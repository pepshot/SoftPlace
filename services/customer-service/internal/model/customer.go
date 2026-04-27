package model

import "github.com/google/uuid"

type Customer struct {
	ID       uuid.UUID
	Login    string
	Email    string
	Password string
}
