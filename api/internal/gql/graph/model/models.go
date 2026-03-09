package model

import "time"

type Recipe struct {
	ID        string
	CreatedAt time.Time

	AuthorID         string
	CurrentVersionID string

	CurrentVersion *RecipeVersion
}

type RecipeVersion struct {
	ID        string
	Recipe    *Recipe
	Version   int32
	Name      string
	PrepMins  int32
	CookMins  int32
	Portions  int32
	CreatedAt time.Time
}
