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
	ContainsGluten     bool
	TaxonomyParentID   *uuid.UUID
	ProcessedLevel     ProcessedLevel
	IsSearchable       bool
}

func NewIngredient(
	name, fileKey string,
	preferredUnit int,
	counter, plural, counterPlural *string,
	animalProductLevel AnimalProductLevel,
	containsGluten bool,
	taxonomyParentID *uuid.UUID,
	processedLevel *int,
	isSearchable *bool,
) (*Ingredient, error) {
	resolvedProcessedLevel := Raw
	if processedLevel != nil {
		resolvedProcessedLevel = ProcessedLevel(*processedLevel)
	}

	resolvedIsSearchable := true
	if isSearchable != nil {
		resolvedIsSearchable = *isSearchable
	}

	if name == "" {
		return nil, ErrInvalidName
	}
	if !isPreferredUnitAllowed(preferredUnit, resolvedIsSearchable) {
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
		ContainsGluten:     containsGluten,
		TaxonomyParentID:   taxonomyParentID,
		ProcessedLevel:     resolvedProcessedLevel,
		IsSearchable:       resolvedIsSearchable,
	}, nil
}

// Unknown unit allowed for non-searchable ingredients
func isPreferredUnitAllowed(preferredUnit int, isSearchable bool) bool {
	resolvedUnit := unit.Unit(preferredUnit)
	if resolvedUnit.IsValid() {
		return true
	}

	return !isSearchable && resolvedUnit == unit.UnitUnknown
}

type AnimalProductLevel int

const (
	Vegan AnimalProductLevel = iota
	Vegetarian
	Meat
)

type ProcessedLevel int

const (
	Unknown ProcessedLevel = iota
	Raw
	Derived
	Composite
)
