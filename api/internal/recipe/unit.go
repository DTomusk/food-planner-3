package recipe

type Unit int

const (
	UnitUnknown Unit = iota
	// A quantum is for when you have x of something, e.g. 5 eggs as opposed to 100ml of eggs
	Quantum
)

func (u Unit) String() string {
	switch u {
	case Quantum:
		return ""
	default:
		return "unknown"
	}
}
