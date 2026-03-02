package reference

type IngredientModel struct {
	ID            string  `yaml:"id"`
	Name          string  `yaml:"name"`
	PreferredUnit int     `yaml:"preferred_unit"`
	Counter       *string `yaml:"counter,omitempty"`
	Plural        *string `yaml:"plural,omitempty"`
	CounterPlural *string `yaml:"counter_plural,omitempty"`
}

type FileData struct {
	Ingredients []*IngredientModel `yaml:"ingredients"`
}
