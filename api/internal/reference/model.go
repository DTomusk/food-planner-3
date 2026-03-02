package reference

type IngredientModel struct {
	ID            string  `yaml:"id"`
	Name          string  `yaml:"name"`
	PreferredUnit int     `yaml:"preferred_unit"`
	Counter       *string `yaml:"counter,omitempty"`
}

type FileData struct {
	Ingredients []*IngredientModel `yaml:"ingredients"`
}
