package tunnel_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/app/tunnel"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestHTTP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start a local HTTP server.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, fmt.Sprintf("You asked for the resource at %s!", r.URL.Path))
		require.NoError(t, err)
	}))
	defer srv.Close()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			tun, err := tunnel.ApplyHTTP(ctx, eit.ControllerClient, client.ObjectKey{Namespace: ns.GetName(), Name: "tunnel"})
			require.NoError(t, err)

			require.NoError(t, tunnel.WithHTTPConnection(ctx, eit.RESTConfig, tun, srv.URL, func(ctx context.Context) {
				script := fmt.Sprintf(`exec wget -qO- %s/foo/bar`, tun.URL())
				r, err := endtoend.Exec(ctx, eit.Environment, script, endtoend.ExecerWithNamespace(ns.GetName()))
				require.NoError(t, err)
				assert.Equal(t, 0, r.Code)
				assert.Equal(t, "You asked for the resource at /foo/bar!", r.Stdout)
				assert.Equal(t, "", r.Stderr)
			}))
		})
	})
}
