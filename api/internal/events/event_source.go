package events

type EventSource string

const (
	AuthServiceSource   EventSource = "auth.service"
	RecipeServiceSource EventSource = "recipe.service"
)
