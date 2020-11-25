package api_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/errawr-go/v2/pkg/testutil"
	"github.com/puppetlabs/leg/httputil/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteErrorWithSensitivity(t *testing.T) {
	tt := []struct {
		Description string
		Context     context.Context
		Error       errawr.Error
	}{
		{
			"Default sensitivity (edge)",
			context.Background(),
			testutil.NewStubError("hello"),
		},
		{
			"Restricted sensitivity (none)",
			api.NewContextWithErrorSensitivity(context.Background(), errawr.ErrorSensitivityNone),
			testutil.NewStubError("hello").WithSensitivity(errawr.ErrorSensitivityEdge),
		},
	}
	for _, test := range tt {
		t.Run(test.Description, func(t *testing.T) {
			w := httptest.NewRecorder()
			api.WriteError(test.Context, w, test.Error)

			var env api.ErrorEnvelope
			require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&env))

			err := env.Error.AsError()
			assert.Equal(t, test.Error.Domain().Key(), err.Domain().Key())
			assert.Equal(t, test.Error.Section().Key(), err.Section().Key())
			assert.Equal(t, test.Error.Sensitivity(), err.Sensitivity())

			guard := errawr.ErrorSensitivityEdge
			if sensitivity, ok := api.ErrorSensitivityFromContext(test.Context); ok {
				guard = sensitivity
			}

			if test.Error.Sensitivity() <= guard {
				assert.Equal(t, test.Error.FormattedDescription().Friendly(), err.FormattedDescription().Friendly())
			} else {
				assert.Equal(t, "", err.FormattedDescription().Friendly())
			}
		})
	}
}
