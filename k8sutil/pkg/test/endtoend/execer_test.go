package endtoend_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestExecer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	script := `
echo foo >&1
echo bar >&2
exit 42
`

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			r, err := endtoend.Exec(ctx, eit.Environment, script, endtoend.ExecerWithNamespace(ns.GetName()))
			require.NoError(t, err)

			assert.Equal(t, 42, r.Code)
			assert.Equal(t, "foo\n", r.Stdout)
			assert.Equal(t, "bar\n", r.Stderr)
		})
	})
}
