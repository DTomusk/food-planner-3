package events

// Maps event types to factory functions
type Registry struct {
	factories map[string]func() Event
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]func() Event),
	}
}

// Adds named event type to registry
// factory returns a zero value type that json unmarshal can populate
func (r *Registry) Register(eventType string, factory func() Event) {
	r.factories[eventType] = factory
}

// Finds the factory for the given event type and calls factory to generate zero value event instance
func (r *Registry) New(eventType string) (Event, bool) {
	factory, ok := r.factories[eventType]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// Returns all registered event type strings
func (r *Registry) Types() []string {
	types := make([]string, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}
