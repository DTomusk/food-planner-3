package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		want      string
	}{
		{
			name:      "returns user agent header value",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			want:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		},
		{
			name:      "trims surrounding whitespace",
			userAgent: "  test-agent/1.0  ",
			want:      "test-agent/1.0",
		},
		{
			name:      "returns empty string when header missing",
			userAgent: "",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}

			require.Equal(t, tt.want, getUserAgent(req))
		})
	}
}

func TestUserAgentMiddleware_AddsUserAgentToContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent/2.0")
	rr := httptest.NewRecorder()

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		userAgent, ok := r.Context().Value(UserAgentKey).(string)
		require.True(t, ok)
		require.Equal(t, "test-agent/2.0", userAgent)
	})

	UserAgentMiddleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
}
