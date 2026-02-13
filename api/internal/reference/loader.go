package reference

type Loader struct {
	filePath string
}

func NewLoader(filePath string) *Loader {
	return &Loader{
		filePath: filePath,
	}
}

func (s *Loader) LoadIngredientData() ([]*IngredientModel, error) {
	return nil, nil
}
