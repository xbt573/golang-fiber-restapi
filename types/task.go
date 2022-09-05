package types

import (
	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID
	Name        string `validate:"required,min=3,max=32"`
	Description string `validate:"max=128"`
}
