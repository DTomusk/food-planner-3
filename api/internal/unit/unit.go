package unit

type Unit int

const (
	UnitUnknown Unit = iota
	// A quantum is for when you have x of something, e.g. 5 eggs as opposed to 100ml of eggs
	Quantum
	Gram
)

func (u Unit) IsValid() bool {
	switch u {
	case Quantum:
		return true
	case Gram:
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
	default:
		return "?"
	}
}
