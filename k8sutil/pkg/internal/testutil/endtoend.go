// Package testutil contains the internal test harness for testing this module.
package testutil

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/test/endtoend"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

type EnvironmentInTest struct {
	*endtoend.Environment
	Labels map[string]string
	t      *testing.T
	nf     endtoend.NamespaceFactory
}

func (eit *EnvironmentInTest) WithNamespace(ctx context.Context, fn func(ns *corev1.Namespace)) {
	require.NoError(eit.t, endtoend.WithNamespace(ctx, eit.Environment, eit.nf, fn))
}

func WithEnvironmentInTest(t *testing.T, fn func(eit *EnvironmentInTest)) {
	viper.SetEnvPrefix("leg_k8sutil_test_e2e")
	viper.AutomaticEnv()

	kubeconfigs := strings.TrimSpace(viper.GetString("kubeconfig"))
	if testing.Short() {
		t.Skip("not running end-to-end tests with -short")
	} else if kubeconfigs == "" {
		t.Skip("not running end-to-end tests without one or more Kubeconfigs specified by LEG_K8SUTIL_TEST_E2E_KUBECONFIG")
	}

	opts := []endtoend.EnvironmentOption{
		endtoend.EnvironmentWithClientKubeconfigs(filepath.SplitList(kubeconfigs)),
		endtoend.EnvironmentWithClientContext(viper.GetString("context")),
	}

	require.NoError(t, endtoend.WithEnvironment(opts, func(e *endtoend.Environment) {
		ls := map[string]string{
			"testutil.internal.k8sutil.leg.puppet.com/harness":   "end-to-end",
			"testutil.internal.k8sutil.leg.puppet.com/test.hash": testHash(t),
		}

		fn(&EnvironmentInTest{
			Environment: e,
			Labels:      ls,
			t:           t,
			nf:          endtoend.NewTestNamespaceFactory(t, endtoend.NamespaceWithLabels(ls)),
		})
	}))
}

func testHash(t *testing.T) string {
	h := sha256.Sum256([]byte(t.Name()))
	return hex.EncodeToString(h[:])[:63]
}
