package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unresolvableObject struct{}

func (unresolvableObject) String() string {
	return "unresolvable"
}

func (unresolvableObject) CacheKey() (string, bool) {
	return "", false
}

type resolvableObject struct {
	value string
}

func (ro resolvableObject) String() string {
	return fmt.Sprintf("resolvable_%s", ro.value)
}

func (ro resolvableObject) CacheKey() (string, bool) {
	return ro.value, true
}

func TestConditionalResolver(t *testing.T) {
	handler := func(object Cacheable) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !NewConditionalResolver(r).Accept(r.Context(), w, object) {
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})
	}

	tests := []struct {
		Name        string
		Object      Cacheable
		Method      string
		IfMatch     string
		IfNoneMatch string
		Expected    int
	}{
		{
			Name:     "resolvable_get_with_matching_etag",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodGet,
			IfMatch:  ETag{Value: "test"}.String(),
			Expected: http.StatusNoContent,
		},
		{
			Name:     "resolvable_get_with_nonmatching_etag",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodGet,
			IfMatch:  ETag{Value: "no"}.String(),
			Expected: http.StatusRequestedRangeNotSatisfiable,
		},
		{
			Name:     "unresolvable_get",
			Object:   &unresolvableObject{},
			Method:   http.MethodGet,
			IfMatch:  ETag{Value: "no"}.String(),
			Expected: http.StatusRequestedRangeNotSatisfiable,
		},
		{
			Name:        "resolvable_get_with_none_match_matching_etag",
			Object:      &resolvableObject{"test"},
			Method:      http.MethodGet,
			IfNoneMatch: ETag{Value: "test"}.String(),
			Expected:    http.StatusNotModified,
		},
		{
			Name:        "resolvable_get_with_none_match_nonmatching_etag",
			Object:      &resolvableObject{"test"},
			Method:      http.MethodGet,
			IfNoneMatch: ETag{Value: "no"}.String(),
			Expected:    http.StatusNoContent,
		},
		{
			Name:     "resolvable_put_with_matching_etag",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodPut,
			IfMatch:  ETag{Value: "test"}.String(),
			Expected: http.StatusNoContent,
		},
		{
			Name:     "resolvable_put_with_nonmatching_etag",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodPut,
			IfMatch:  ETag{Value: "no"}.String(),
			Expected: http.StatusPreconditionFailed,
		},
		{
			Name:     "unresolvable_put",
			Object:   &unresolvableObject{},
			Method:   http.MethodPut,
			IfMatch:  ETag{Value: "no"}.String(),
			Expected: http.StatusPreconditionFailed,
		},
		{
			Name:        "resolvable_put_with_none_match_matching_etag",
			Object:      &resolvableObject{"test"},
			Method:      http.MethodPut,
			IfNoneMatch: ETag{Value: "test"}.String(),
			Expected:    http.StatusPreconditionFailed,
		},
		{
			Name:        "resolvable_put_with_none_match_nonmatching_etag",
			Object:      &resolvableObject{"test"},
			Method:      http.MethodPut,
			IfNoneMatch: ETag{Value: "no"}.String(),
			Expected:    http.StatusNoContent,
		},
		{
			Name:        "unresolvable_put_with_none_match",
			Object:      &unresolvableObject{},
			Method:      http.MethodPut,
			IfNoneMatch: ETag{Value: "no"}.String(),
			Expected:    http.StatusNoContent,
		},
		{
			Name:     "resolvable_delete_with_no_etag",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodDelete,
			Expected: http.StatusNoContent,
		},
		{
			Name:     "unresolvable_delete_with_no_etag",
			Object:   &unresolvableObject{},
			Method:   http.MethodDelete,
			Expected: http.StatusNoContent,
		},
		{
			Name:     "resolvable_put_with_match_any",
			Object:   &resolvableObject{"test"},
			Method:   http.MethodPut,
			IfMatch:  "*",
			Expected: http.StatusNoContent,
		},
		{
			Name:     "unresolvable_put_with_match_any",
			Object:   &unresolvableObject{},
			Method:   http.MethodPut,
			IfMatch:  "*",
			Expected: http.StatusNoContent,
		},
		{
			Name:     "nil_put_with_match_any",
			Object:   nil,
			Method:   http.MethodPut,
			IfMatch:  "*",
			Expected: http.StatusPreconditionFailed,
		},
		{
			Name:        "resolvable_put_with_none_match_any",
			Object:      &resolvableObject{"test"},
			Method:      http.MethodPut,
			IfNoneMatch: "*",
			Expected:    http.StatusPreconditionFailed,
		},
		{
			Name:        "unresolvable_put_with_none_match_any",
			Object:      &unresolvableObject{},
			Method:      http.MethodPut,
			IfNoneMatch: "*",
			Expected:    http.StatusPreconditionFailed,
		},
		{
			Name:        "nil_put_with_none_match_any",
			Object:      nil,
			Method:      http.MethodPut,
			IfNoneMatch: "*",
			Expected:    http.StatusNoContent,
		},
		{
			Name:     "nil_put",
			Object:   nil,
			Method:   http.MethodPut,
			IfMatch:  ETag{Value: "test"}.String(),
			Expected: http.StatusPreconditionFailed,
		},
		{
			Name:        "nil_put_with_none_match",
			Object:      nil,
			Method:      http.MethodPut,
			IfNoneMatch: ETag{Value: "test"}.String(),
			Expected:    http.StatusNoContent,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req := httptest.NewRequest(test.Method, "/", nil)
			if test.IfMatch != "" {
				req.Header.Set("if-match", test.IfMatch)
			} else if test.IfNoneMatch != "" {
				req.Header.Set("if-none-match", test.IfNoneMatch)
			}

			res := httptest.NewRecorder()
			handler(test.Object).ServeHTTP(res, req)

			assert.Equal(t, test.Expected, res.Code)
		})
	}
}

func TestScanETags(t *testing.T) {
	tests := []struct {
		Header   string
		OK       bool
		Expected []ETag
	}{
		{
			Header: `"foo", "bar", W/"baz"`,
			OK:     true,
			Expected: []ETag{
				{Value: "foo"},
				{Value: "bar"},
				{Weak: true, Value: "baz"},
			},
		},
		{
			Header: `invalid`,
			OK:     false,
		},
		{
			Header:   `*`,
			OK:       true,
			Expected: []ETag{},
		},
	}
	for _, test := range tests {
		tags, ok := scanETags([]string{test.Header})

		assert.Equal(t, test.OK, ok, "for header %q", test.Header)
		assert.Equal(t, test.Expected, tags, "for header %q", test.Header)
	}
}
