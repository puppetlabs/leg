package selfsignedsecret_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/app/selfsignedsecret"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/testutil"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestReconcilerSetsCertificateInSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testutil.WithEnvironmentInTest(t, func(eit *testutil.EnvironmentInTest) {
		eit.WithNamespace(ctx, func(ns *corev1.Namespace) {
			secretKey := client.ObjectKey{
				Namespace: ns.GetName(),
				Name:      "test",
			}

			mgr, err := manager.New(eit.RESTConfig, manager.Options{
				Namespace: ns.GetName(),
			})
			require.NoError(t, err)

			selfsignedsecret.AddReconcilerToManager(mgr, secretKey, "Puppet", "sss.example.com")

			var wg sync.WaitGroup
			defer wg.Wait()

			mctx, cancel := context.WithCancel(ctx)
			defer cancel()

			wg.Add(1)
			go func() {
				defer wg.Done()
				require.NoError(t, mgr.Start(mctx))
			}()

			secret := corev1obj.NewTLSSecret(secretKey)
			require.NoError(t, secret.Persist(ctx, eit.ControllerClient))

			require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
				if ok, err := secret.Load(ctx, eit.ControllerClient); err != nil {
					return true, err
				} else if !ok {
					return true, fmt.Errorf("secret disappeared")
				}

				if len(secret.Object.Data["tls.key"]) == 0 {
					return false, fmt.Errorf("secret has not had a certificate added")
				}

				return true, nil
			}))
		})
	})
}
