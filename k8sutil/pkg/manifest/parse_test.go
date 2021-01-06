package manifest_test

import (
	"os"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestParse(t *testing.T) {
	f, err := os.Open("testdata/fixtures/ambassador-v1.8.0-rbac.yaml")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, f.Close())
	}()

	objs, err := manifest.Parse(scheme.Scheme, f)
	require.NoError(t, err)
	require.Len(t, objs, 5)
	assert.IsType(t, &corev1.Service{}, objs[0])
	assert.IsType(t, &rbacv1beta1.ClusterRole{}, objs[1])
	assert.IsType(t, &corev1.ServiceAccount{}, objs[2])
	assert.IsType(t, &rbacv1beta1.ClusterRoleBinding{}, objs[3])
	assert.IsType(t, &appsv1.Deployment{}, objs[4])
}
