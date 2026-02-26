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
