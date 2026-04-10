package ingredient

import (
	"foodplanner/internal/unit"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID                 uuid.UUID
	Name               string
	FileKey            string
	PreferredUnit      unit.Unit
	Counter            *string
	Plural             *string
	CounterPlural      *string
	AnimalProductLevel AnimalProductLevel
}

func NewIngredient(name, fileKey string, preferredUnit int, counter, plural, counterPlural *string, animalProductLevel AnimalProductLevel) (*Ingredient, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if !unit.Unit(preferredUnit).IsValid() {
		return nil, ErrInvalidPreferredUnit
	}
	return &Ingredient{
		ID:                 uuid.New(),
		Name:               name,
		FileKey:            fileKey,
		PreferredUnit:      unit.Unit(preferredUnit),
		Counter:            counter,
		Plural:             plural,
		CounterPlural:      counterPlural,
		AnimalProductLevel: animalProductLevel,
	}, nil
}

type AnimalProductLevel int

const (
	Vegan AnimalProductLevel = iota
	Vegetarian
	Meat
)
