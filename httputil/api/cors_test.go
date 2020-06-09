package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCORSBuilder(t *testing.T) {
	handler := NewCORSBuilder().
		AllowOrigins("http://example.com", "http://app.example.com").
		AllowHeaderPrefix("horsehead-").
		AllowHeaders("X-Custom-Header").Build()

	req, err := http.NewRequest(http.MethodOptions, "http://example.com", nil)
	require.NoError(t, err)
	req.Header.Set("access-control-request-method", "POST")
	req.Header.Set("access-control-request-headers", "Horsehead-Custom-Header, X-Custom-Header")
	req.Header.Set("Origin", "http://app.example.com")

	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	result := resp.Result()

	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, "http://app.example.com", result.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "Horsehead-Custom-Header, X-Custom-Header", result.Header.Get("Access-Control-Allow-Headers"))
	require.Equal(t, strings.Join(corsDefaultAllowedMethods, ", "), result.Header.Get("Access-Control-Allow-Methods"))

	{
		// a request that should fail
		req, err := http.NewRequest(http.MethodOptions, "http://example.com", nil)
		require.NoError(t, err)
		req.Header.Set("access-control-request-method", "PUT")

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)
		result := resp.Result()

		require.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	}
}
