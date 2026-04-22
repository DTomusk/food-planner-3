package reference

type IngredientModel struct {
	ID                 string  `yaml:"id"`
	Name               string  `yaml:"name"`
	PreferredUnit      int     `yaml:"preferred_unit"`
	Counter            *string `yaml:"counter,omitempty"`
	Plural             *string `yaml:"plural,omitempty"`
	CounterPlural      *string `yaml:"counter_plural,omitempty"`
	AnimalProductLevel int     `yaml:"animal_product_level,omitempty"`
	ContainsGluten     bool    `yaml:"contains_gluten,omitempty"`
	TaxonomyParentKey  *string `yaml:"taxonomy_parent_key,omitempty"`
	ProcessedLevel     *int    `yaml:"processed_level,omitempty"`
	IsSearchable       *bool   `yaml:"is_searchable,omitempty"`
}

type FileData struct {
	Ingredients []*IngredientModel `yaml:"ingredients"`
}
