package recipe

import "github.com/google/uuid"

type Recipe struct {
	ID   uuid.UUID
	Name string
}
