package reference

type IngredientModel struct {
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	PreferredUnit int    `yaml:"preferred_unit"`
	FileKey       string `yaml:"file_key"`
}

type FileData struct {
	Ingredients []*IngredientModel `yaml:"ingredients"`
}
