package ingredient

import "errors"

var (
	ErrInvalidName          = errors.New("invalid ingredient name")
	ErrInvalidPreferredUnit = errors.New("invalid preferred unit")
)
