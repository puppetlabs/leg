package jsonpath_test

import (
	"context"
	"os"
	"testing"

	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type regressionQuery struct {
	ID       string         `yaml:"id"`
	Selector string         `yaml:"selector"`
	Document any            `yaml:"document"`
	Ordered  *bool          `yaml:"ordered,omitempty"`
	Results  map[string]any `yaml:",inline"`
}

type regressionSuite struct {
	Queries []*regressionQuery `yaml:"queries"`
}

// TestRegression runs all of the tests in a regression test suite that follows
// the format of the files in the repository
// https://github.com/cburgmer/json-path-comparison. Because that repository is
// licensed under the GPL, none of its contents are included here (and will not
// be accepted in the future).
//
// Use the LEG_JSONUTIL_JSONPATH_REGRESSION_SUITE_PATH environment variable to
// point to such a file on your filesystem.
//
// Note that you should not expect this test to pass. Make informed decisions
// about what you want to improve from them, but as there's no actual JSONPath
// standard, if our behavior exists as a superset of what's expected by the
// user, we're still pretty much fine.
func TestRegression(t *testing.T) {
	p := os.Getenv("LEG_JSONUTIL_JSONPATH_REGRESSION_SUITE_PATH")
	if p == "" {
		t.Skip("not running regression suite without LEG_JSONUTIL_JSONPATH_REGRESSION_SUITE_PATH")
	}

	fp, err := os.Open(p)
	require.NoError(t, err)
	defer fp.Close()

	var suite regressionSuite
	require.NoError(t, yaml.NewDecoder(fp).Decode(&suite))

	for _, q := range suite.Queries {
		t.Run(q.ID, func(t *testing.T) {
			t.Logf("document: %+v", q.Document)
			t.Logf("selector: %q", q.Selector)

			consensus, found := q.Results["consensus"]
			if !found {
				t.Log("no consensus")
				return
			}

			ctx := context.Background()

			get, err := jsonpath.New(q.Selector)
			if consensus == "NOT_SUPPORTED" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			value, err := get(ctx, q.Document)
			if q.Results["not-found-consensus"] == "NOT_FOUND" && (err != nil || value == nil) {
				// OK, although we may also match on the consensus.
				return
			}
			require.NoError(t, err)

			if scalarConsensus, found := q.Results["scalar-consensus"]; found {
				require.Equal(t, scalarConsensus, value)
			} else if q.Ordered != nil && !*q.Ordered {
				require.ElementsMatch(t, value, consensus)
			} else {
				require.Equal(t, consensus, value)
			}
		})
	}
}
