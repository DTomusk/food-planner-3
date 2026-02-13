package reference

type ReferenceService struct{}

func NewReferenceService() *ReferenceService {
	return &ReferenceService{}
}

func (s *ReferenceService) LoadIngredientData() ([]*IngredientModel, error) {
	return nil, nil
}
