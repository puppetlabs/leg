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
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestHTTPS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start a local HTTP server.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, fmt.Sprintf("You asked for the resource at %s via %s!", r.URL.Path, r.Host))
		require.NoError(t, err)
	}))
	defer srv.Close()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			tun, err := tunnel.ApplyHTTPS(ctx, eit.ControllerClient, client.ObjectKey{Namespace: ns.GetName(), Name: "tunnel"})
			require.NoError(t, err)

			cert, err := tun.CertificateAuthorityPEM()
			require.NoError(t, err)

			require.NoError(t, tunnel.WithHTTPConnection(ctx, eit.RESTConfig, tun.HTTP, srv.URL, func(ctx context.Context) {
				// Wait for service.
				_, err = corev1obj.NewEndpointsBoundPoller(corev1obj.NewEndpoints(tun.TLSProxy.Service)).Load(ctx, eit.ControllerClient)
				require.NoError(t, err)

				script := fmt.Sprintf(`
mkdir -p /etc/ssl/certs
cat >>/etc/ssl/certs/ca-certificates.crt <<'EOT'
%s
EOT
exec wget -qO- %s/foo/bar
`, cert, tun.URL())
				r, err := endtoend.Exec(ctx, eit.Environment, script, endtoend.ExecerWithNamespace(ns.GetName()))
				require.NoError(t, err)
				assert.Equal(t, 0, r.Code)
				assert.Equal(t, fmt.Sprintf("You asked for the resource at /foo/bar via %s!", tun.TLSProxy.Service.DNSName()), r.Stdout)
				assert.Equal(t, "", r.Stderr)
			}))
		})
	})
}
