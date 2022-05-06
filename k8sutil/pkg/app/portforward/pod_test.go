package portforward_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/app/portforward"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestForwardPod(t *testing.T) {
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

			// Forward to pod and get response.
			require.NoError(t, portforward.ForwardPod(ctx, eit.RESTConfig, pod, []uint16{8080}, func(ctx context.Context, m portforward.Map) error {
				require.NotEqual(t, 0, m[8080])
				require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
					resp, err := http.Get(fmt.Sprintf("http://localhost:%d", m[8080]))
					if err != nil {
						return retry.Repeat(err)
					}
					defer resp.Body.Close()

					b, err := ioutil.ReadAll(resp.Body)
					require.NoError(t, err)
					require.Equal(t, []byte("Hello, world!\n"), b)

					return retry.Done(nil)
				}))

				return nil
			}))
		})
	})
}
