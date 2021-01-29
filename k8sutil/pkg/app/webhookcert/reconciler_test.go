package webhookcert_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/app/selfsignedsecret"
	"github.com/puppetlabs/leg/k8sutil/pkg/app/webhookcert"
	admissionregistrationv1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/admissionregistrationv1"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestReconcilerUpdatesCABundle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			secret := corev1obj.NewTLSSecret(client.ObjectKey{
				Namespace: ns.GetName(),
				Name:      "test",
			})
			vwc := admissionregistrationv1obj.NewValidatingWebhookConfiguration(ns.GetName())
			vwc.Object.Webhooks = []admissionregistrationv1.ValidatingWebhook{
				{
					Name: "test.webhookcert.example.com",
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Namespace: ns.GetName(),
							Name:      "test",
						},
					},
					SideEffects:             func(s admissionregistrationv1.SideEffectClass) *admissionregistrationv1.SideEffectClass { return &s }(admissionregistrationv1.SideEffectClassNone),
					AdmissionReviewVersions: []string{"v1"},
				},
			}
			defer func() {
				_, err := vwc.Delete(ctx, eit.ControllerClient)
				require.NoError(t, err)
			}()

			mgr, err := manager.New(eit.RESTConfig, manager.Options{
				Namespace: ns.GetName(),
			})
			require.NoError(t, err)

			selfsignedsecret.AddReconcilerToManager(mgr, secret.Key, "Puppet", "test.webhookcert.example.com")
			webhookcert.AddReconcilerToManager(mgr, secret.Key, webhookcert.WithValidatingWebhookConfiguration(vwc.Name))

			var wg sync.WaitGroup
			defer wg.Wait()

			mctx, cancel := context.WithCancel(ctx)
			defer cancel()

			wg.Add(1)
			go func() {
				defer wg.Done()
				require.NoError(t, mgr.Start(mctx))
			}()

			require.NoError(t, secret.Persist(ctx, eit.ControllerClient))
			require.NoError(t, vwc.Persist(ctx, eit.ControllerClient))

			require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
				if ok, err := vwc.Load(ctx, eit.ControllerClient); err != nil {
					return true, err
				} else if !ok {
					return true, fmt.Errorf("configuration disappeared")
				}

				if len(vwc.Object.Webhooks[0].ClientConfig.CABundle) == 0 {
					return false, fmt.Errorf("configuration does not contain a CA bundle")
				}

				return true, nil
			}))
		})
	})
}
