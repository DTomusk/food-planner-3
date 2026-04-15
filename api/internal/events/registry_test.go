package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_RegisterAndNew_ReturnsFactoryInstance(t *testing.T) {
	registry := NewRegistry()
	registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

	event, ok := registry.New(UserSignedUpType)
	require.True(t, ok)
	require.IsType(t, &UserSignedUpEvent{}, event)
}

func TestRegistry_New_UnknownType_ReturnsFalse(t *testing.T) {
	registry := NewRegistry()

	event, ok := registry.New("missing.type")
	require.False(t, ok)
	require.Nil(t, event)
}

func TestRegistry_Types_ContainsRegisteredTypes(t *testing.T) {
	registry := NewRegistry()
	registry.Register("type.a", func() Event { return &UserSignedUpEvent{} })
	registry.Register("type.b", func() Event { return &UserSignedUpEvent{} })

	types := registry.Types()
	require.ElementsMatch(t, []string{"type.a", "type.b"}, types)
}
