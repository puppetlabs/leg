package tlsproxy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/app/tlsproxy"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestTLSProxy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			// Create a pod to serve test output.
			pod := corev1obj.NewPod(client.ObjectKey{Namespace: ns.GetName(), Name: "test"})
			pod.Object.Spec = corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "server",
						Image: "hashicorp/http-echo",
						Args: []string{
							"-listen", ":8080",
							"-text", "Hello, world!",
						},
						Ports: []corev1.ContainerPort{
							{
								Name:          "http",
								ContainerPort: 8080,
								Protocol:      corev1.ProtocolTCP,
							},
						},
					},
				},
			}
			require.NoError(t, pod.Persist(ctx, eit.ControllerClient))

			// Wait for pod to start.
			_, err := corev1obj.NewPodRunningPoller(pod).Load(ctx, eit.ControllerClient)
			require.NoError(t, err)

			// Apply TLS proxy for pod.
			tp, err := tlsproxy.Apply(
				ctx,
				eit.ControllerClient,
				client.ObjectKey{Namespace: ns.GetName(), Name: "proxy"},
				fmt.Sprintf("%s:8080", pod.Object.Status.PodIP),
			)
			require.NoError(t, err)

			// Wait for service.
			_, err = corev1obj.NewEndpointsBoundPoller(corev1obj.NewEndpoints(tp.Service)).Load(ctx, eit.ControllerClient)
			require.NoError(t, err)

			// Check to make sure our proxy works!
			cert, err := tp.CertificateAuthorityPEM()
			require.NoError(t, err)

			script := fmt.Sprintf(`
mkdir -p /etc/ssl/certs
cat >>/etc/ssl/certs/ca-certificates.crt <<'EOT'
%s
EOT
exec wget -qO- %s
`, cert, tp.URL())
			r, err := endtoend.Exec(ctx, eit.Environment, script, endtoend.ExecerWithNamespace(ns.GetName()))
			require.NoError(t, err)
			assert.Equal(t, 0, r.Code)
			assert.Equal(t, "Hello, world!\n", r.Stdout)
			assert.Equal(t, "", r.Stderr)
		})
	})
}
