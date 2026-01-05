package recipe

import "testing"

func TestInstantiateRecipe(t *testing.T) {
	name := "Pancakes"
	recipe, err := NewRecipe(name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if recipe.Name != name {
		t.Errorf("Expected name %q, got %q", name, recipe.Name)
	}
}

func TestEmptyRecipeName(t *testing.T) {
	_, err := NewRecipe("")
	if err != ErrEmptyName {
		t.Fatalf("Expected error %v, got %v", ErrEmptyName, err)
	}
}
