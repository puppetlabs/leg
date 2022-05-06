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
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestForwardService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			// Create a pod to serve test output.
			pod := corev1obj.NewPod(client.ObjectKey{Namespace: ns.GetName(), Name: "test-0"})
			pod.Object.SetLabels(eit.Labels)
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

			// Wrap pod with service.
			svc := corev1obj.NewService(client.ObjectKey{Namespace: ns.GetName(), Name: "test"})
			svc.Object.Spec = corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
				Ports: []corev1.ServicePort{
					{
						Name:       "http",
						TargetPort: intstr.FromString("http"),
						Protocol:   corev1.ProtocolTCP,
						Port:       80,
					},
				},
				Selector: pod.Object.GetLabels(),
			}
			require.NoError(t, svc.Persist(ctx, eit.ControllerClient))

			// Forward to service and get response.
			require.NoError(t, portforward.ForwardService(ctx, eit.RESTConfig, svc, 80, func(ctx context.Context, port uint16) error {
				require.NotEqual(t, 0, port)
				require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
					resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
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
