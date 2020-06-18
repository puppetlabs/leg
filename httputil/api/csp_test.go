package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSPBuilder(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	builder := CSPBuilder{}
	builder.SetDirective(CSPDefaultSrc, "self").
		SetDirective(CSPScriptSrc, "self", "*.example.com").
		SetDirective(CSPMediaSrc, "none").
		SetDirective(CSPImgSrc, "example-bucket-gGi7b2.s3.amazonaws.com").
		SetDirective(CSPBlockAllMixedContent)

	wrapped := builder.Middleware(handler)

	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()

	wrapped.ServeHTTP(resp, req)
	result := resp.Result()

	require.Equal(t, http.StatusOK, result.StatusCode)

	csp := result.Header.Get("content-security-policy")
	require.Equal(t, "default-src 'self'; script-src 'self' *.example.com; media-src 'none'; img-src example-bucket-gGi7b2.s3.amazonaws.com; block-all-mixed-content", csp)
}
