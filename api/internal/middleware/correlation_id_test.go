package middleware

import (
	"foodplanner/internal/correlationid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCorrelationIDMiddleware_GeneratesCorrelationIDWhenHeaderMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		id := correlationid.FromContext(r.Context())
		require.NotEqual(t, uuid.Nil, id)
	})

	CorrelationIDMiddleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
}

func TestCorrelationIDMiddleware_UsesHeaderCorrelationIDWhenValid(t *testing.T) {
	expectedID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Correlation-ID", expectedID.String())
	rr := httptest.NewRecorder()

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		id := correlationid.FromContext(r.Context())
		require.Equal(t, expectedID, id)
	})

	CorrelationIDMiddleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
}

func TestCorrelationIDMiddleware_IgnoresInvalidHeaderCorrelationID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Correlation-ID", "not-a-uuid")
	rr := httptest.NewRecorder()

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		id := correlationid.FromContext(r.Context())
		require.NotEqual(t, uuid.Nil, id)
	})

	CorrelationIDMiddleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
}
