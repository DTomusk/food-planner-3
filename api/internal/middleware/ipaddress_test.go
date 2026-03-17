package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIPAddress(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name: "prefers first forwarded IP",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.10, 198.51.100.20",
				"X-Real-IP":       "192.0.2.30",
			},
			remoteAddr: "127.0.0.1:8080",
			want:       "203.0.113.10",
		},
		{
			name: "uses real IP when forwarded header is absent",
			headers: map[string]string{
				"X-Real-IP": "192.0.2.30",
			},
			remoteAddr: "127.0.0.1:8080",
			want:       "192.0.2.30",
		},
		{
			name:       "uses remote address host when proxy headers are absent",
			remoteAddr: "198.51.100.25:9090",
			want:       "198.51.100.25",
		},
		{
			name:       "falls back to raw remote address when host port parsing fails",
			remoteAddr: "malformed-remote-addr",
			want:       "malformed-remote-addr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			require.Equal(t, tt.want, getIPAddress(req))
		})
	}
}

func TestIPMiddleware_AddsIPAddressToContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 198.51.100.20")
	rr := httptest.NewRecorder()

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		ip, ok := r.Context().Value(IPKey).(string)
		require.True(t, ok)
		require.Equal(t, "203.0.113.10", ip)
	})

	IPMiddleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
}
