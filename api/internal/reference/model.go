package reference

type IngredientModel struct {
	FileKey            string             `yaml:"id"`
	Name               string             `yaml:"name"`
	PreferredUnit      int                `yaml:"preferred_unit"`
	Counter            *string            `yaml:"counter,omitempty"`
	Plural             *string            `yaml:"plural,omitempty"`
	CounterPlural      *string            `yaml:"counter_plural,omitempty"`
	AnimalProductLevel int                `yaml:"animal_product_level,omitempty"`
	ContainsGluten     bool               `yaml:"contains_gluten,omitempty"`
	ProcessedLevel     *int               `yaml:"processed_level,omitempty"`
	IsSearchable       *bool              `yaml:"is_searchable,omitempty"`
	Children           []*IngredientModel `yaml:"children,omitempty"`
}

type FileData struct {
	Ingredients []*IngredientModel `yaml:"ingredients"`
}
