package recipe

type RecipeSource struct {
	Type         SourceType
	URL          *string
	BookTitle    *string
	BookPage     *int32
	Instructions *string
}

type SourceType int

const (
	SourceUnknown SourceType = iota
	URL
	BookReference
	Original
)

func newSource(sourceRequest *CreateRecipeSourceRequest) (*RecipeSource, error) {
	if sourceRequest == nil {
		return nil, ErrNoSource
	}
	switch sourceRequest.Type {
	case 0:
		return nil, nil
	case 1:
		if sourceRequest.URL == nil {
			return nil, ErrMissingURL
		}
		return NewURLSource(*sourceRequest.URL)
	case 2:
		if sourceRequest.BookTitle == nil || sourceRequest.BookPage == nil {
			return nil, ErrMissingBookReference
		}
		return NewBookReferenceSource(*sourceRequest.BookTitle, *sourceRequest.BookPage)
	case 3:
		if sourceRequest.Instructions == nil {
			return nil, ErrMissingInstructions
		}
		return NewOriginalSource(*sourceRequest.Instructions)
	default:
		return nil, ErrInvalidSourceType
	}
}

// TODO: add stricter validity for URLs
func NewURLSource(url string) (*RecipeSource, error) {
	if url == "" {
		return nil, ErrEmptyURL
	}
	return &RecipeSource{
		Type: URL,
		URL:  &url,
	}, nil
}

func NewBookReferenceSource(bookTitle string, bookPage int32) (*RecipeSource, error) {
	if bookTitle == "" {
		return nil, ErrEmptyBookTitle
	}
	if bookPage <= 0 {
		return nil, ErrInvalidBookPage
	}
	return &RecipeSource{
		Type:      BookReference,
		BookTitle: &bookTitle,
		BookPage:  &bookPage,
	}, nil
}

func NewOriginalSource(instructions string) (*RecipeSource, error) {
	if instructions == "" {
		return nil, ErrEmptyInstructions
	}
	return &RecipeSource{
		Type:         Original,
		Instructions: &instructions,
	}, nil
}
