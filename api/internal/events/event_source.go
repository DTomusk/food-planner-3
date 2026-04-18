package events

type EventSource string

const (
	AuthServiceSource   EventSource = "auth.service"
	GraphQLServerSource EventSource = "graphql.server"
	RecipeServiceSource EventSource = "recipe.service"
)
