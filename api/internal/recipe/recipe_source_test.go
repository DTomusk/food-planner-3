package recipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewURLSource(t *testing.T) {
	// Arrange
	url := "https://example.com/recipe"

	// Act
	source, err := NewURLSource(url)

	// Assert
	require.NoError(t, err)
	require.Equal(t, URL, source.Type)
	require.Equal(t, url, *source.URL)
}

func TestNewURLSource_EmptyURL(t *testing.T) {
	// Act
	source, err := NewURLSource("")

	// Assert
	require.Error(t, err)
	require.Nil(t, source)
	require.Equal(t, ErrEmptyURL, err)
}

func TestNewBookReferenceSource(t *testing.T) {
	// Arrange
	bookTitle := "Pee pee poo poo book"
	bookPage := int32(4)

	// Act
	source, err := NewBookReferenceSource(bookTitle, bookPage)

	// Assert
	require.NoError(t, err)
	require.Equal(t, BookReference, source.Type)
	require.Equal(t, bookTitle, *source.BookTitle)
	require.Equal(t, bookPage, *source.BookPage)
}

func TestNewBookReferenceSource_EmptyBookTitle(t *testing.T) {
	// Act
	source, err := NewBookReferenceSource("", int32(4))

	// Assert
	require.Error(t, err)
	require.Nil(t, source)
	require.Equal(t, ErrEmptyBookTitle, err)
}

func TestNewBookReferenceSource_InvalidBookPage(t *testing.T) {
	// Act
	source, err := NewBookReferenceSource("Pee pee poo poo book", int32(0))

	// Assert
	require.Error(t, err)
	require.Nil(t, source)
	require.Equal(t, ErrInvalidBookPage, err)
}

func TestNewOriginalSource(t *testing.T) {
	// Arrange
	instructions := "Make the dish"

	// Act
	source, err := NewOriginalSource(instructions)

	// Assert
	require.NoError(t, err)
	require.Equal(t, Original, source.Type)
	require.Equal(t, instructions, *source.Instructions)
}

func TestNewOriginalSource_EmptyInstructions(t *testing.T) {
	// Act
	source, err := NewOriginalSource("")

	// Assert
	require.Error(t, err)
	require.Nil(t, source)
	require.Equal(t, ErrEmptyInstructions, err)
}
