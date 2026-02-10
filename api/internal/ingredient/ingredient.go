package ingredient

import "github.com/google/uuid"

type Ingredient struct {
	ID   uuid.UUID
	Name string
}
