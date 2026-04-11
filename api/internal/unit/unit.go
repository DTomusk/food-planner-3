package unit

type Unit int

const (
	UnitUnknown Unit = iota
	// A quantum is for when you have x of something, e.g. 5 eggs as opposed to 100ml of eggs
	Quantum
	Gram
	Milliliter
	Tablespoon
	Teaspoon
	Centimeter
)

func (u Unit) IsValid() bool {
	switch u {
	case Quantum:
		return true
	case Gram:
		return true
	case Milliliter:
		return true
	case Tablespoon:
		return true
	case Teaspoon:
		return true
	case Centimeter:
		return true
	default:
		return false
	}
}

func (u Unit) String() string {
	switch u {
	case Quantum:
		return "Count"
	case Gram:
		return "Grams"
	case Milliliter:
		return "Milliliters"
	case Tablespoon:
		return "Tablespoons"
	case Teaspoon:
		return "Teaspoons"
	case Centimeter:
		return "Centimeters"
	default:
		return "Unknown"
	}
}

func (u Unit) Symbol() string {
	switch u {
	case Quantum:
		return ""
	case Gram:
		return "g"
	case Milliliter:
		return "ml"
	case Tablespoon:
		return "tbsp"
	case Teaspoon:
		return "tsp"
	case Centimeter:
		return "cm"
	default:
		return "?"
	}
}
